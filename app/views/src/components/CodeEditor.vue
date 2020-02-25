<template>
  <div class="vue-cm">
    <textarea ref="cm-area"></textarea>
  </div>
</template>

<script>
import CodeMirror from 'codemirror'
import 'codemirror/lib/codemirror.css'

const CmJs = () => import('codemirror/mode/javascript/javascript')
const CmCl = () => import('codemirror/mode/clike/clike')
const CmPy = () => import('codemirror/mode/python/python')
const CmXml = () => import('codemirror/mode/xml/xml')
const CmVue = () => import('codemirror/mode/vue/vue')
const CmCss = () => import('codemirror/mode/css/css')
const CmMd = () => import('codemirror/mode/markdown/markdown')
const CmSh = () => import('codemirror/mode/shell/shell')
const CmLua = () => import('codemirror/mode/lua/lua')

const tmMaterial = () => import('codemirror/theme/material.css')
const tmMaterialDark = () => import('codemirror/theme/material-darker.css')
const tmMaterialOcean = () => import('codemirror/theme/material-ocean.css')
const tmMaterialPalenight = () => import('codemirror/theme/material-palenight.css')

export default {
  name: 'codedisplay',
  props: {
    value: {
      type: String,
      default: ''
    },
    lang: {
      type: String,
      default: null
    },
    height: {
      type: String,
      default: 'auto'
    },
    theme: {
      type: String,
      default: 'default'
    },
    options: {
      type: Object,
      default: () => ({
        tabSize: 2,
        lineNumbers: true,
        line: true
      })
    }
  },
  data: () => ({
    cminstance: null,
    modes: [
      {
        value: 'css',
        alias: 'Css',
        script: CmCss
      }, {
        value: 'javascript',
        alias: 'Javascript',
        script: CmJs
      }, {
        value: 'html',
        alias: 'XML/HTML',
        script: CmXml
      }, {
        value: 'x-python',
        alias: 'Python',
        script: CmPy
      }, {
        value: 'x-vue',
        alias: 'Vue',
        script: CmVue
      }, {
        value: 'markdown',
        alias: 'Markdown',
        script: CmMd
      }, {
        value: 'x-objectivec',
        alias: 'C',
        script: CmCl
      }, {
        value: 'x-sh',
        alias: 'Shell',
        script: CmSh
      }, {
        value: 'x-lua',
        alias: 'lua',
        script: CmLua
      }
    ],
    themes: [
      {
        name: 'material',
        full: 'material',
        loader: tmMaterial
      }, {
        name: 'darker',
        full: 'material-darker',
        loader: tmMaterialDark
      }, {
        name: 'ocean',
        full: 'material-ocean',
        loader: tmMaterialOcean
      }, {
        name: 'palenight',
        full: 'material-palenight',
        loader: tmMaterialPalenight
      }, {
        name: 'default',
        full: 'default',
        loader: null
      }
    ]
  }),
  mounted () {
    this.initialize()
  },
  watch: {
    value: function (newval, oldval) {
      if (this.cminstance.getValue() !== newval) {
        this.cminstance.setValue(newval)
      }
    },
    lang: function (newval, oldval) {
      if (oldval !== newval) {
        this.loadMode()
      }
    },
    height: function (newval, oldval) {
      if (oldval !== newval) {
        this.cminstance.setSize(null, newval)
      }
    },
    theme: function (newval, oldval) {
      if (oldval !== newval) {
        this.loadTheme()
      }
    },
    options: function (newval, oldval) {
      if (oldval !== newval) {
        this.loadOptions()
      }
    }
  },
  methods: {
    initialize () {
      this.cminstance = CodeMirror.fromTextArea(this.$refs['cm-area'])
      this.cminstance.getWrapperElement().style.fontFamily = 'SFMono-Regular, Consolas, Liberation Mono, Menlo, monospace'
      this.cminstance.getWrapperElement().style.borderRadius = '5px'
      this.cminstance.refresh()
      this.cminstance.setValue(this.value)
      this.cminstance.setSize(null, this.height)
      this.cminstance.on('change', (cminstance) => {
        const value = cminstance.getValue()
        if (this.$emit && this.value !== value) {
          this.$emit('input', value)
          // cconsole.log(cminstance.getValue())
          this.$emit('change')
        }
      })
      this.loadMode()
      this.loadOptions()
      this.loadTheme()
    },
    loadMode () {
      let given = ''
      if (this.lang) {
        given = this.lang
      }
      this.cminstance.setOption('mode', 'text')
      this.modes.forEach(mode => {
        if (mode.alias.toLowerCase() === given.toLowerCase() ||
            mode.value.toLowerCase() === given.toLowerCase()) {
          mode.script().then(() => {
            this.cminstance.setOption('mode', 'text/' + mode.value)
          })
        }
      })
    },
    loadOptions () {
      // theme and mode are not set here
      this.options.mode = undefined
      this.options.theme = undefined
      for (const key in this.options) {
        if (this.options[key] !== undefined) {
          this.cminstance.setOption(key, this.options[key])
        }
      }
    },
    loadTheme () {
      let given = 'default'
      if (this.theme) {
        given = this.theme
      }
      this.themes.forEach(theme => {
        if (theme.name.toLowerCase() === given.toLowerCase() ||
            theme.full.toLowerCase() === given.toLowerCase()) {
          if (theme.loader) {
            theme.loader().then(() => {
              this.cminstance.setOption('theme', theme.full)
            })
          }
        }
      })
    }
  }
}
</script>
