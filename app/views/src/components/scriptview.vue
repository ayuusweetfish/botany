<template>
  <div>
    <el-card shadow="never" style="margin: 120px 0px 20px 0px">
      <div align="left" style="font-size: 18px; font-weight: 600; margin-bottom: 20px">执行脚本</div>
      <el-row style="margin-bottom: 20px">
        <el-col :span="4">
          <div align="right" style="margin-top: 8px; margin-right: 20px">输入参数</div>
        </el-col>
        <el-col :span="18">
          <el-input size="small" v-model="input">
            <el-button slot="append" @click="passArg">运行</el-button>
          </el-input>
        </el-col>
      </el-row>
    </el-card>
    <el-card shadow="never" v-loading="logLoading">
      <div align="left" style="font-size: 18px; font-weight: 600; margin-bottom: 10px">脚本日志</div>
      <div align="left" style="margin-bottom: 10px">
        <el-link type="primary" :underline="false" icon="el-icon-download" @click="downLoadLog">下载全部日志</el-link>
      </div>
      <div class="cm-container">
        <codemirror v-model="log" :options="cmOptions" align="left"></codemirror>
      </div>
    </el-card>
  </div>
</template>

<script>
import { codemirror } from 'vue-codemirror-lite'

export default {
  name: 'contestscript',
  components: {
    codemirror
  },
  created () {
    this.cid = this.$route.query.cid
    this.getLog()
  },
  data () {
    return {
      input: '',
      log: '',
      cid: '',
      logLoading: false,
      cmOptions: {
        lineNumbers: true,
        lineWrapping: true,
        indentUnit: 2,
        tabSize: 2,
        autoCloseBrackets: true,
        readOnly: true
      }
    }
  },
  methods: {
    passArg () {
      const loading = this.$loading({ lock: true, text: '处理中' })
      const params = this.$qs.stringify({
        arg: this.input
      })
      this.$axios.post(
        '/contest/' + this.cid + '/manual_script',
        params
      ).then(res => {
        loading.close()
        this.$message.success('操作成功')
        this.getLog()
      }).catch(err => {
        console.log(err)
        loading.close()
        this.$message.error('操作失败')
      })
    },
    getLog () {
      this.logLoading = true
      this.$axios.get(
        '/contest/' + this.cid + '/script_log',
        { params: { full: 0 } }
      ).then(res => {
        this.log = res.data
        this.logLoading = false
      }).catch(err => {
        console.log(err)
        this.$message.error('查询失败')
        this.logLoading = false
      })
    },
    downLoadLog () {
      const loading = this.$loading({ lock: true, text: '请求中' })
      this.$axios.get(
        '/contest/' + this.cid + '/script_log',
        { params: { full: 1 } }
      ).then(res => {
        loading.close()
        const blob = new Blob([res.data], { type: 'text/txt,charset=UTF-8' })
        const a = document.createElement('a')
        a.href = URL.createObjectURL(blob)
        a.download = 'log.txt'
        a.click()
      }).catch(err => {
        console.log(err)
        loading.close()
        this.$message.error('请求失败')
      })
    }
  }
}
</script>

<style>
.cm-container {
  border: 1px solid #dcdfe6;
  max-height: 600px;
}
</style>
