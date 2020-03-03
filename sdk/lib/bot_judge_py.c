#define PY_SSIZE_T_CLEAN
#include <Python.h>
#include "bot.h"

static PyObject *init(PyObject *self, PyObject *args)
{
    PyObject *list;
    if (!PyArg_ParseTuple(args, "O", &list)) return NULL;

    if (!PyList_Check(list)) return NULL;
    PyObject *iter = PyObject_GetIter(list);
    // assert(iter != NULL);

    Py_ssize_t len = PySequence_Length(list);
    char **argv = (char **)malloc(len * sizeof(char *));

    Py_ssize_t i;
    for (i = 0; i < len; i++) {
        PyObject *elm = PyIter_Next(iter);
        // assert(elm != NULL);

        if (!PyUnicode_Check(elm)) return NULL;

        PyObject *tmp = PyUnicode_AsEncodedString(elm, "UTF-8", "strict");
        if (!tmp) return NULL;

        const char *str = PyBytes_AS_STRING(tmp);
        argv[i] = strdup(str);

        Py_DECREF(tmp);
    }

    int n = bot_judge_init(len, argv);

    for (i = 0; i < len; i++) free(argv[i]);
    free(argv);

    return PyLong_FromLong(n);
}

static PyObject *send(PyObject *self, PyObject *args)
{
    int id;
    const char *str;
    if (!PyArg_ParseTuple(args, "is", &id, &str)) return NULL;

    bot_judge_send(id, str);

    return Py_None;
}

static PyObject *recv(PyObject *self, PyObject *args)
{
    int id;
    int timeout;
    if (!PyArg_ParseTuple(args, "ii", &id, &timeout)) return NULL;

    int len;
    char *resp = bot_judge_recv(id, &len, timeout);

    PyObject *ret = Py_BuildValue("si", resp, (int)len);
    return ret;
}

static PyObject *finish(PyObject *self, PyObject *args)
{
    bot_judge_finish();
    return Py_None;
}

static PyMethodDef methods[] = {
    {"init", init, METH_VARARGS, NULL},
    {"send", send, METH_VARARGS, NULL},
    {"recv", recv, METH_VARARGS, NULL},
    {"finish", finish, METH_VARARGS, NULL},
    {NULL, NULL, 0, NULL}
};

static struct PyModuleDef module = {
    PyModuleDef_HEAD_INIT,
    "bot_judge_py",
    "Botany judge module (Python binding)",
    -1,
    methods
};

PyMODINIT_FUNC PyInit_bot_judge_py()
{
    return PyModule_Create(&module);
}
