#!/bin/bash

tool=$1
arg1=$2
arg2=$3

_script_name_="tools.sh"
_script_call_help_str_=". ${_script_name_}"

_is_number_regexp_='^[0-9]+$'

# Вызов без аргументов -- открыть справку.
if [[ ${tool} == '' && ${arg1} == '' ]]; then
  tool="help"
fi

# ----------------------------------------------------------------------------------------------------------------------
# HELP
# ----------------------------------------------------------------------------------------------------------------------

if [[ ${tool} == 'help' ]]; then
  help_tool=${arg1}

  # HELP
  if [[ ${help_tool} == '' || ${help_tool} == 'help' ]]; then
    echo ""
    echo "Шаблон:"
    echo "    ${_script_call_help_str_} TOOL [ARG_1 [ARG_2 ...]]]"
    echo ""
    echo "Аргумент TOOL:"
    echo "    help    -- эта справка."
    echo "    up      -- развернуть все."
    echo "    down    -- свернуть все."
    echo "    migrate -- для работы с миграциями."
    echo ""
    echo "Аргументы ARG_...:"
    echo "    Зависят от конкретной тулы -- см. \"${_script_call_help_str_} help TOOL\""
    echo ""

  # UP
  elif [[ ${help_tool} == 'up' ]]; then
    echo ""
    echo "Шаблон:"
    echo "    ${_script_call_help_str_} up [dev [--no-migrate]]"
    echo ""
    echo "Аргумент dev:"
    echo "    dev -- в режиме разработки. --no-migrate -- пропустить запуск миграций"
    echo ""

  # DOWN
  elif [[ ${help_tool} == 'down' ]]; then
    echo ""
    echo "Шаблон:"
    echo "    ${_script_call_help_str_} down"
    echo ""

  # MIGRATE
  elif [[ ${help_tool} == 'migrate' ]]; then
    echo ""
    echo "Шаблон:"
    echo "    ${_script_call_help_str_} migrate ACTION"
    echo ""
    echo "Аргументы:"
    echo "    ACTION -- основаны на \"example/migrator\" контейнере:"
    docker build -q --tag="example/migrator:latest" .migrator/build > /dev/null
    docker run --rm example/migrator:latest help
    echo ""

  # ОШИБКА
  else
    echo ""
    echo "Нет справки для тулы \"${help_tool}\"!"
    echo ""
    return
  fi

# ----------------------------------------------------------------------------------------------------------------------
# UP
# ----------------------------------------------------------------------------------------------------------------------

elif [[ ${tool} == 'up' ]]; then
  up_mode=${arg1}
  up_opt1=${arg2}

  # DEV
  if [[ ${up_mode} == 'dev' ]]; then
    if [[ ! (${up_opt1} == "" || ${up_opt1} == "--no-migrate") ]]; then
      echo "Неизвестная опция: ${up_opt1}"
      return
    fi

    docker compose rm --stop --force
    docker compose build

    if [[ ${up_opt1} == "--no-migrate" ]]; then
      echo "Запуск миграций пропущен!"
    else
      source ${_script_name_} migrate up
    fi

    # Watch при ребилдах создает "безымянные" (sha256:...) образы, которые накапливаются и не удаляются автоматически.
    # Это быстро приводит к заполнению диска, и приходится прерывать процесс, чтобы "down" все почистил.
    # Так что пусть лучше чистит все в параллельном фоновом процессе.
    echo -n "Запускается фоновый уборщик watch-мусора..."
    bg_watch_gb_collector_cmd="yes | docker system prune --volumes > /dev/null 2>&1"
    nohup watch -n 120 "$bg_watch_gb_collector_cmd" > /dev/null 2>&1 &
    nohup_pid=$!
    for _ in {1..10}; do
        sleep 1
        echo -n "." # без переноса строки
        watch_prune_bg_pid=$(pgrep -f "$bg_watch_gb_collector_cmd")
        if [[ $watch_prune_bg_pid ]]; then
            echo "" # перенос строки
            break
        fi
    done
    if [[ $watch_prune_bg_pid ]]; then
      echo "Фоновый уборщик watch-мусора - $nohup_pid / $watch_prune_bg_pid"
    else
      echo "Фоновый уборщик watch-мусора не запустился :("
      kill $nohup_pid || true
      return
    fi

    docker compose up --remove-orphans --watch --menu
    # Тут "docker compose up --watch" блокирует терминал -- ждем выхода.
    # --

    echo "Киллится фоновый уборщик watch-мусора..."
    kill $watch_prune_bg_pid 2>/dev/null || echo "уже"
    kill $nohup_pid 2>/dev/null || echo "уже"
    echo "Фоновый уборщик watch-мусора килльнут!"

    # Хоть сейчас "безымянные" образы и удаляются в фоне, не помешает выполнить очистку и после завершения процесса.
    # Вызов здесь "docker image (или system) prune" подошел бы, но т.к. это dev, вызов тулы "down" даже лучше.
    source ${_script_name_} down
    # При выходе через Ctrl+Z "docker compose up --watch" все еще существует -- повторные "up dev" накопят его дубликаты.
    # Поэтому нужно завершать все "docker compose up", запущенные из этой (чтоб не прибить другие проекты) директории.
    for word in $(ps -x | grep "docker compose up")
    do
      if [[ ${word} =~ ${_is_number_regexp_} ]] ; then
          pid=${word}
          if [[ -d /proc/${pid} && $(readlink /proc/${pid}/cwd) == $(pwd) ]]; then
            kill -9 ${pid}
          fi
      fi
    done

  # PROD
  elif [[ ${up_mode} == "" ]]; then
    docker compose up --remove-orphans --build --detach
    source ${_script_name_} migrate up
    docker system prune --force

  # ОШИБКА
  else
    echo ""
    echo "Неизвестный аргумент \"${up_mode}\"!"
    echo ""
    return
  fi

# ----------------------------------------------------------------------------------------------------------------------
# DOWN
# ----------------------------------------------------------------------------------------------------------------------

elif [[ ${tool} == 'down' ]]; then
  docker compose --profile=script-containers down   # по профилю + все без профиля (обычные сервисы).
  docker system prune --force

# ----------------------------------------------------------------------------------------------------------------------
# MIGRATE
# ----------------------------------------------------------------------------------------------------------------------

elif [[ ${tool} == 'migrate' ]]; then
  migrate_action=${arg1}

  if [[ ${migrate_action} == '' ]]; then
    echo ""
    echo "Надо указать ACTION!"
    echo ""
    return
  fi

  migration_service_name="script-migrate-${migrate_action}"

  is_service_exist=$(docker compose --profile=script-containers config --format=json | grep "${migration_service_name}")
  if [[ ! ${is_service_exist} ]]; then
    echo ""
    echo "Нет такого сервиса \"${migration_service_name}\""
    echo ""
    return
  fi

  # Если у скрипта есть зависимости (например, у migrate up, вероятно, какие-то базы данных),
  # то их контейнеры не будут остановлены после завершения работы контейнера скрипта.
  # Логично их останавливать, если они не были запущены ранее.

  # Собираем ранее незапущенные зависимости
  # --------------------------------
  python_script=""
  python_script="${python_script}import json;"
  python_script="${python_script}print("
  python_script="${python_script}' '.join("
  python_script="${python_script}json.loads('''$(docker compose --profile=script-containers config --format=json)''')"
  python_script="${python_script}['services']['${migration_service_name}'].get('depends_on',{}).keys()"
  python_script="${python_script}))"
  migration_service_dependencies=$(docker run --rm --entrypoint=python python:3-alpine -c "${python_script}")

  launched_services=$(docker compose ps --services)
  built_services=$(docker compose ps --services --all)
  previously_not_launched_dependencies=""
  previously_not_built_dependencies=""
  for dependency_service in ${migration_service_dependencies}
  do
    is_launched=0
    for launched_service in ${launched_services}
    do
      if [[ "${dependency_service}" == "${launched_service}" ]]; then
        is_launched=1
        break
      fi
    done

    is_built=0
    for built_service in ${built_services}
    do
      if [[ "${dependency_service}" == "${built_service}" ]]; then
        is_built=1
        break
      fi
    done

    if [[ ${is_launched} == 0 ]]; then
      previously_not_launched_dependencies="${previously_not_launched_dependencies} ${dependency_service}"
    fi
    if [[ ${is_built} == 0 ]]; then
      previously_not_built_dependencies="${previously_not_built_dependencies} ${dependency_service}"
    fi
  done

  # Выполнение команды
  # --------------------------------

  docker compose build ${migration_service_name}
  docker compose --profile=script-containers run --rm ${migration_service_name}

  # Остановка ранее незапущенных зависимостей
  # --------------------------------

  launched_services=$(docker compose ps --services)
  for launched_service in ${launched_services}
  do
    for service_to_stop in ${previously_not_launched_dependencies}
    do
      if [[ "${launched_service}" == "${service_to_stop}" ]]; then
        docker compose stop ${service_to_stop}
        break
      fi
    done

    for service_to_remove in ${previously_not_built_dependencies}
    do
      if [[ "${launched_service}" == "${service_to_remove}" ]]; then
        docker compose rm --force ${service_to_remove}
        break
      fi
    done
  done

# ----------------------------------------------------------------------------------------------------------------------
# ОШИБКА
# ----------------------------------------------------------------------------------------------------------------------

else
  echo "Нет такой тулы \"${tool}\"!"
fi
