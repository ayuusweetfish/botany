<template>
  <div>
    <el-card style="margin: 10px 0px 20px 0px">
      <div align="left" style="font-size: 18px; font-weight: 600; margin-bottom: 20px">执行脚本</div>
      <el-row>
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
    <el-card v-loading="logLoading">
      <div align="left" style="font-size: 18px; font-weight: 600; margin-bottom: 20px">脚本日志</div>
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
      const loading = this.$loading({lock: true, text: '处理中'})
      let params = this.$qs.stringify({
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
        '/contest/' + this.cid + '/script_log'
      ).then(res => {
        this.log = res.data
        this.logLoading = false
      }).catch(err => {
        console.log(err)
        this.$message.error('查询失败')
        this.logLoading = false
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
