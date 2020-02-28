<template>
  <v-container fluid>
    <v-snackbar
      top
      style="margin-top: 60px"
      v-model="error.show"
      color="error"
      :timeout="3000">
      {{error.msg}}
    </v-snackbar>
    <v-snackbar
      top
      style="margin-top: 60px"
      v-model="success.show"
      color="success"
      :timeout="3000">
      {{success.msg}}
    </v-snackbar>
    <v-row justify="center">
      <v-col :cols="12" :md="10" :lg="8">
        <v-card outlined style="border: none">
          <v-card-title class="title">脚本操作</v-card-title>
          <v-card-subtitle class="subtitle-1">当前比赛ID：{{$route.params.cid}}，名称：{{$store.state.cname}}</v-card-subtitle>
          <v-card-text>
            <div class="mt-2 mb-2 body-1">运行脚本</div>
            <div class="d-flex">
              <v-text-field
                v-model="command"
                outlined dense
                label="输入命令"
              ></v-text-field>
              <v-btn
                color="primary"
                height="40"
                :loading="submiting"
                :disabled="submiting"
                @click="submit"
              ><v-icon>mdi-play</v-icon>执行</v-btn>
            </div>
            <div>
              <a
                :href="`${$axios.defaults.baseURL}/contest/${$route.params.cid}/script_log?full=1`"
                download="log.txt"
                type="text/txt,charset=UTF-8"
                class="body-1"
                style="text-decoration: none"
              ><v-icon color="primary">mdi-download-outline</v-icon>下载全部赛事日志</a>
            </div>
          </v-card-text>
        </v-card>
        <v-card outlined>
          <code-editor
            v-model="log"
            theme="material"
            :options="cmOptions"
            height="auto"
          ></code-editor>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import CodeEditor from '../components/CodeEditor.vue'
export default {
  name: 'HandleScript',
  components: {
    'code-editor': CodeEditor
  },
  mounted () {
    this.getLog()
  },
  data: () => ({
    command: '',
    log: '',
    loading: false,
    submiting: false,
    success: {
      msg: '',
      show: false
    },
    error: {
      msg: '',
      show: false
    },
    cmOptions: {
      lineNumbers: true,
      indentUnit: 2,
      tabSize: 2,
      autoCloseBrackets: true,
      readOnly: true
    }
  }),
  methods: {
    submit () {
      this.submiting = true
      const params = this.$qs.stringify({
        arg: this.input
      })
      this.$axios.post(
        '/contest/' + this.$route.params.cid + '/manual_script',
        params
      ).then(res => {
        this.submiting = false
        this.success.msg = '提交成功'
        this.success.show = true
        this.getLog()
      }).catch(() => {
        this.submiting = false
        this.error.msg = '提交失败'
        this.error.show = true
      })
    },
    getLog () {
      this.loading = true
      this.$axios.get(
        '/contest/' + this.$route.params.cid + '/script_log',
        { params: { full: 0 } }
      ).then(res => {
        this.log = res.data
        this.loading = false
      }).catch(() => {
        this.loading = false
        this.error.msg = '查询日志失败'
        this.error.show = true
      })
    }
  }
}
</script>

<style>

</style>
