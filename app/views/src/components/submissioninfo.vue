<template>
  <div>
    <el-row style="margin-bottom: 10px">
      <el-card shadow="never">
        <div align="left">
          <div>提交时间：{{topbarTime}}</div>
          <div>提交编号：{{topbarID}}</div>
          <div>代码语言：{{lang}}</div>
          <div v-if="topbarInfo.status">
            <div style="display: inline">代码状态：</div><div :style="'display: inline; color: ' + topbarColor.status">{{topbarInfo.status}}</div>
          </div>
          <div v-if="topbarInfo.msg">
            <div style="display: inline">编译信息：</div><div :style="'display: inline; color: ' + topbarColor.msg">{{topbarInfo.msg}}</div>
          </div>
        </div>
      </el-card>
    </el-row>
    <el-row style="margin-bottom: 10px">
      <el-card body-style="min-height: 360px" v-loading="codeLoading"  shadow="never">
        <div align="left" style="font-size: 18px; font-weight: 600; margin-bottom: 20px">代码查看</div>
        <div class="cm-container">
          <codemirror v-model="code" :options="cmOptions" align="left"></codemirror>
        </div>
      </el-card>
    </el-row>
  </div>
</template>

<script>
import { codemirror } from 'vue-codemirror-lite'

export default {
  name: 'submissioninfo',
  components: {
    codemirror
  },
  created () {
    this.cid = this.$route.query.cid
    this.sid = this.$route.query.sid
    this.showCode()
  },
  data () {
    return {
      code: '',
      lang: '',
      codeLoading: false,
      topbarTime: '',
      topbarID: '',
      topbarInfo: {
        status: '',
        msg: ''
      },
      topbarColor: {
        status: '',
        msg: ''
      },
      sid: '',
      cid: '',
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
    showCode () {
      this.codeLoading = true
      this.$axios.get(
        '/contest/' + this.cid + '/submission/' + this.sid
      ).then(res => {
        this.codeLoading = false
        this.code = res.data.contents
        this.lang = res.data.lang
        this.topbarTime = this.$functions.dateTimeString(res.data.created_at)
        this.topbarID = this.sid
        const statcolor = this.getStatColor(res.data.status)
        this.topbarInfo.status = statcolor.stat
        this.topbarColor.status = statcolor.color
        this.topbarColor.msg = 'gray'
        this.topbarInfo.msg = res.data.msg
      }).catch(err => {
        this.codeLoading = false
        console.log(err)
      })
    },
    getStatColor (statid) {
      let color = ''
      let stat = ''
      switch (statid) {
        case this.$consts.codeStat.pending:
          stat = '等待处理'
          color = 'grey'
          break
        case this.$consts.codeStat.compiling:
          stat = '编译中'
          color = 'orange'
          break
        case this.$consts.codeStat.accepted:
          stat = '接受'
          color = 'green'
          break
        case this.$consts.codeStat.complErr:
          stat = '编译错误'
          color = 'red'
          break
        case this.$consts.codeStat.systmErr:
          stat = '系统错误(请联系管理员)'
          color = 'red'
          break
        default:
          break
      }
      return {
        color: color,
        stat: stat
      }
    }
  }
}
</script>

<style scoped>
.cm-container {
  border: 1px solid #dcdfe6;
  max-height: 600px;
}
</style>
