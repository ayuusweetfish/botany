<template>
  <div>
    <el-row style="margin-bottom: 10px">
      <el-card>
        <div align="left">
          <div>{{topbarText}}</div>
          <div v-if="topbarInfo.status">
            <div style="display: inline">代码状态：</div>
            <div :style="'display: inline; color: ' + topbarColor.status">{{topbarInfo.status}}</div>
          </div>
          <div v-if="topbarInfo.msg">
            <div style="display: inline">编译信息：</div>
            <div :style="'display: inline; color: ' + topbarColor.msg">{{topbarInfo.msg}}</div>
          </div>
        </div>
      </el-card>
    </el-row>
    <el-row :gutter="20">
      <el-col :span="18">
        <el-row style="margin-bottom: 10px">
          <el-card body-style="min-height: 360px">
            <div align="left" style="font-size: 18px; font-weight: 600; margin-bottom: 20px">代码编辑</div>
            <div class="cm-container">
              <codemirror v-model="code" :options="cmOptions" align="left"></codemirror>
            </div>
          </el-card>
        </el-row>
        <el-row>
          <el-card>
            <div align="left" style="font-size: 18px; font-weight: 600; margin-bottom: 20px">操作</div>
            <el-row>
              <el-col :span="12">
                <el-button type="primary" size="small" style="width: 80%" @click="submitCode">提交</el-button>
              </el-col>
              <el-col :span="12">
                <el-button size="small" style="width: 80%" @click="clearCode">清空</el-button>
              </el-col>
            </el-row>
          </el-card>
        </el-row>
      </el-col>
      <el-col :span="6">
        <el-card body-style="min-height: 480px">
          <div align="left" style="font-size: 18px; font-weight: 600; margin-bottom: 20px">历史代码</div>
          <el-timeline align="left">
            <el-timeline-item
              v-for="(activity, index) in history"
              :key="index"
              :timestamp="activity.time"
              :color="activity.color"
              placement="top"
            >
              <div align="left">编号：{{activity.sid}}</div>
              <div align="left">状态：{{activity.stat}}</div>
              <el-button type="text" size="small" @click="showCode(activity.sid)">点击导出</el-button>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { codemirror } from 'vue-codemirror-lite'

export default {
  name: 'submissions',
  components: {
    codemirror
  },
  created () {
    this.cid = this.$route.query.id
    this.getHistory()
  },
  data () {
    return {
      code: '',
      history: [],
      topbarText: '尚未提交代码',
      topbarInfo: {
        status: '',
        msg: ''
      },
      topbarColor: {
        status: '',
        msg: ''
      },
      cid: '',
      cmOptions: {
        lineNumbers: true,
        lineWrapping: true,

        indentUnit: 2,
        tabSize: 2,
        autoCloseBrackets: true
      }
    }
  },
  methods: {
    getHistory () {
      const loading = this.$loading({lock: true, text: '加载中'})
      this.$axios.get(
        '/contest/' + this.cid + '/my'
      ).then(res => {
        console.log(res.data)
        this.history = []
        res.data.forEach(item => {
          let statcolor = this.getStatColor(item.status)
          this.history.push({
            sid: item.id,
            time: this.$functions.dateTimeString(item.created_at),
            color: statcolor.color,
            stat: statcolor.stat
          })
        })
        if (this.history.length > 0) {
          this.topbarText = '代码已提交，点击"导出"以查看'
        } else {
          this.topbarText = '尚未提交代码'
        }
        loading.close()
      }).catch(err => {
        loading.close()
        if (err.response.status === 403 || err.response.status === 401) {
          this.topbarText = '尚未登录或未参加比赛'
        } else {
          this.$message.error('查询失败')
        }
      })
    },
    clearCode () {
      this.code = ''
    },
    submitCode () {
      const loading = this.$loading({lock: true, text: '上传中'})
      let params = this.$qs.stringify({
        code: this.code
      })
      this.$axios.post(
        '/contest/' + this.cid + '/submit',
        params
      ).then(res => {
        loading.close()
        this.$message.success('提交成功')
        window.location.reload()
      }).catch(err => {
        loading.close()
        console.log(err)
      })
    },
    showCode (sid) {
      const loading = this.$loading({lock: true, text: '查询中'})
      this.$axios.get(
        '/contest/' + this.cid + '/submission/' + sid
      ).then(res => {
        loading.close()
        this.code = res.data.contents
        this.topbarText = '提交于' + this.$functions.dateTimeString(res.data.created_at) + '，编号为' + sid
        let statcolor = this.getStatColor(res.data.status)
        this.topbarInfo.status = statcolor.stat
        this.topbarColor.status = statcolor.color
        this.topbarColor.msg = 'gray'
        this.topbarInfo.msg = res.data.msg
      }).catch(err => {
        loading.close()
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

<style>
.cm-container .CodeMirror {
  height: auto;
  max-height: 600px;
  font-family: monospace;
  position: relative;
  background: white;
  direction: ltr;
  overflow: hidden;
}
.cm-container .CodeMirror-wrap {
  height: auto;
  font-family: monospace;
  position: relative;
  background: white;
  direction: ltr;
}
.cm-container .CodeMirror-scroll {
  min-height: 300px;
  max-height: 600px;
}
</style>
