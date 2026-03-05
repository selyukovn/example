import importlib
import os
import pydoc
import sys

sys.path.insert(0,'/migrations/packages')

DIR_SCRIPTS = 'scripts'
DIR_STATE = 'state'


def log(msg):
    print('%s' % msg)

# ----------------

processed_script_names = [
    name
    for name in os.listdir(DIR_STATE)
    if (
        name[0] == 's'
        and name.split('.').__len__() == 1
    )
]

current_script_names = [
    filename[:-3]
    for filename in os.listdir(DIR_SCRIPTS)
    if (
        filename[0] == 's'
        and filename[-3:] == '.py'
    )
]

# ----------------

script_names_to_execute = [
    name
    for name in current_script_names
    if name not in processed_script_names
]

if script_names_to_execute.__len__() == 0:
    log('Нечего выполнять!')
    exit(0)

script_names_to_execute.sort()

log('К выполнению (%s):\r\n%s' % (
    script_names_to_execute.__len__(),
    '\r\n'.join(script_names_to_execute)
))

# ----------------

for script_name in script_names_to_execute:
    log('%s: Загрузка...' % script_name)

    importlib.__import__('%s.%s' % (DIR_SCRIPTS, script_name))
    script = pydoc.locate('%s.%s.%s' % (DIR_SCRIPTS, script_name, script_name))()

    log('%s: Выполнение...' % script_name)

    script.up()

    log('%s: Сохранение состояния...' % script_name)

    state_file = open(os.path.join(DIR_STATE, script_name), 'w')
    state_file.close()

    log('%s: Готово!' % script_name)

log('Готово!')

exit(0)
