
cmd=$1

if [[ ${cmd} == '' || ${cmd} == 'help' ]]; then
  echo '"up" -- накатывает ВСЕ еще не примененные миграции;'
  echo '"down" -- откатывает ПОСЛЕДНЮЮ примененную миграцию;'
  echo '"new" -- создает новый скрипт в папке-источнике миграций;'
  echo '"help" -- отображает эту информацию;'
  return
fi

if [[ ${cmd} == 'up' || ${cmd} == 'down' ]]; then
  echo 'Запуск "'${cmd}'"...'
  echo 'Выполнение...'
  python ${cmd}.py
  echo 'Готово!'

  # При выполнении скриптов создается и лезет в git.
  rm -Rf ./scripts/__pycache__

  return
fi

if [[ ${cmd} == 'new' ]]; then
  script_name=s$(date +"%Y%m%d%H%M%S")
  echo 'Создание '${script_name}'.py ...'

  echo 'class '${script_name}':'  > scripts/${script_name}.py
  echo '    def up(self):'        >> scripts/${script_name}.py
  echo '        pass'             >> scripts/${script_name}.py
  echo '    def down(self):'      >> scripts/${script_name}.py
  echo '        pass'             >> scripts/${script_name}.py

  echo 'Готово!'
  return
fi

echo 'Неизвестная команда "'${cmd}'"!'
