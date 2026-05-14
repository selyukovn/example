import importlib
import os
import pydoc
import sys

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

if processed_script_names.__len__() == 0:
    log('Нечего выполнять!')
    exit(0)

processed_script_names.sort()
script_name = processed_script_names[-1]

# ----------------

log('%s: Загрузка...' % script_name)

importlib.__import__('%s.%s' % (DIR_SCRIPTS, script_name))
script = pydoc.locate('%s.%s.%s' % (DIR_SCRIPTS, script_name, script_name))()

log('%s: Выполнение...' % script_name)

script.down()

log('%s: Сохранение состояния...' % script_name)

os.remove(os.path.join(DIR_STATE, script_name))

log('%s: Готово!' % script_name)
