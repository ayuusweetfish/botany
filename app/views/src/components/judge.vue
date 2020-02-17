<template>
  <div>
    <el-row v-loading="loading" style="margin-bottom: 10px">
      <el-card shadow="never">
        <div align="left">
          <div v-if="topbarID">
            <div>提交时间：{{topbarTime}}</div>
            <div>提交编号：{{topbarID}}</div>
          </div>
          <div v-else>{{topbarText}}</div>
          <div>代码语言：<el-select
              v-model="lang"
              size="small"
            >
              <el-option v-for="(item, key) in langSel" :key="key" :value="item">{{item}}</el-option>
            </el-select>
          </div>
          <div v-if="topbarInfo.status">
            <div style="display: inline">代码状态：</div><div :style="'display: inline; color: ' + topbarColor.status">{{topbarInfo.status}}</div>
          </div>
          <div v-if="topbarInfo.msg">
            <div style="display: inline">编译信息：</div><div :style="'display: inline; color: ' + topbarColor.msg">{{topbarInfo.msg}}</div>
          </div>
        </div>
      </el-card>
    </el-row>
    <el-row :gutter="16">
      <el-col :span="17">
        <el-row style="margin-bottom: 10px">
          <div style="min-height: 360px" v-loading="codeLoading" shadow="never">
            <div class="cm-container">
              <codemirror v-model="code" :options="cmOptions" align="left"></codemirror>
            </div>
          </div>
        </el-row>
        <el-row>
          <div>
            <el-row>
              <el-col :span="12" align="left">
                <el-button type="primary" size="small" style="width: 80%" @click="submitCode" :disabled="!canSubmit">{{submitText}}</el-button>
              </el-col>
              <el-col :span="12" align="right">
                <el-button size="small" style="width: 80%" @click="clearCode">清空</el-button>
              </el-col>
            </el-row>
          </div>
        </el-row>
      </el-col>
      <el-col :span="7">
        <div style="min-height: 480px; margin-left: 20px;">
          <div align="left" style="font-size: 18px; font-weight: 600; margin-bottom: 20px">裁判代码</div>
            <div v-if="mainCode.sid!==''">
              <div align="left">编号：{{mainCode.sid}}</div>
              <div align="left">语言：{{mainCode.lang}}</div>
              <el-button type="text" size="small" @click="showCode(mainCode.sid)">导出</el-button>
              <el-button type="text" size="small" @click="setJudge(-1)">撤下</el-button>
            </div>
            <div v-else style="margin-bottom: 20px; font-size: 14px; color: gray">你还没有设置裁判程序</div>
          <div align="left" style="font-size: 18px; font-weight: 600; margin-bottom: 20px">其他历史代码</div>
          <el-timeline align="left">
            <el-timeline-item
              v-for="(activity, index) in historyPart"
              :key="index"
              :timestamp="activity.time"
              :color="activity.color"
              reverse
              placement="top"
            >
              <div align="left">编号：{{activity.sid}}</div>
              <div align="left">语言：{{activity.lang}}</div>
              <div align="left">状态：{{activity.stat}}</div>
              <el-button type="text" size="small" @click="showCode(activity.sid)">导出</el-button>
              <el-button
                v-if="mainCode.sid !== activity.sid"
                type="text"
                size="small"
                @click="setJudge(activity.sid)"
                :disabled="activity.stat!=='接受'"
              >设为裁判</el-button>
              <el-button
                v-else
                type="text" size="small" @click="showCode(activity.sid)"
                disabled
              >裁判代码</el-button>
            </el-timeline-item>
          </el-timeline>
          <el-pagination
            :total="history.length"
            :current-page="page"
            :page-size="3"
            @current-change="handleCurrentChange"
            layout="prev, pager, next"
            size="small"
          >
          </el-pagination>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { codemirror } from 'vue-codemirror-lite'

export default {
  name: 'judge',
  components: {
    codemirror
  },
  created () {
    this.cid = this.$route.query.cid
    this.getHistory()
  },
  data () {
    return {
      loading: false,
      code: '',
      lang: '',
      page: 1,
      mainCode: {
        sid: '',
        time: '',
        lang: '',
        color: '',
        stat: ''
      },
      codeLoading: false,
      canSubmit: true,
      submitText: '提交(作为新记录)',
      history: [],
      historyPart: [],
      topbarTime: '',
      topbarID: '',
      topbarText: '新代码',
      topbarInfo: {
        status: '',
        msg: ''
      },
      topbarColor: {
        status: '',
        msg: ''
      },
      langSel: ['gcc.c', 'gcc.cpp', '5.1.lua', '3.py'],
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
    handleCurrentChange (val) {
      this.page = val
      this.updateHistoryPart()
    },
    updateHistoryPart () {
      const start = (this.page - 1) * 3
      const end = start + 3
      this.historyPart = this.history.slice(start, end)
    },
    // compareTime () {
    //   this.$axios.get(
    //     '/contest/' + this.cid + '/info'
    //   ).then(res => {
    //     const time1 = res.data.start_time
    //     const time2 = res.data.end_time
    //     const now = new Date().getTime() / 1000
    //     if (time1 < now && now < time2) {
    //       this.canSubmit = true
    //       this.submitText = '提交(作为新记录)'
    //     } else if (time1 > now) {
    //       this.canSubmit = false
    //       this.submitText = '还不能提交'
    //     } else {
    //       this.canSubmit = false
    //       this.submitText = '已超过提交时间'
    //     }
    //   }).catch(err => {
    //     console.log(err)
    //   })
    // },
    setJudge (sid) {
      this.loading = true
      const params = this.$qs.stringify({
        submission: sid
      })
      this.$axios.post(
        '/contest/' + this.cid + '/judge',
        params
      ).then(res => {
        this.loading = false
        this.$message.success('设置成功')
        this.getHistory()
      }).catch(err => {
        this.loading = false
        console.log(err)
      })
    },
    getHistory () {
      this.loading = true
      // this.compareTime()
      this.mainCode = {
        sid: '',
        time: '',
        lang: '',
        color: '',
        stat: ''
      }
      this.$axios.get(
        '/contest/' + this.cid + '/my'
      ).then(res => {
        console.log(res.data)
        this.history = []
        this.$axios.get(
          '/contest/' + this.cid + '/judge_id'
        ).then(result => {
          console.log(result.data)
          if (result.data.submission === -1) {
            this.mainCode = {
              sid: '',
              time: '',
              lang: '',
              color: '',
              stat: ''
            }
          }
          res.data.forEach(item => {
            const statcolor = this.getStatColor(item.status)
            const historyItem = {
              sid: item.id,
              time: this.$functions.dateTimeString(item.created_at),
              color: statcolor.color,
              stat: statcolor.stat,
              lang: item.lang
            }
            if (historyItem.sid === result.data.judge) {
              this.mainCode = historyItem
            }
            if (historyItem.sid !== -1) {
              this.history.push(historyItem)
            }
          })
          if (this.history.length === 0 && this.mainCode.sid === '') {
            this.topbarText = '尚未提交代码'
          }
          this.page = 1
          this.updateHistoryPart()
          this.loading = false
        }).catch(err => {
          console.log(err)
            .loading = false
        })
      }).catch(err => {
        this.loading = false
        if (err.response.status === 403 || err.response.status === 401) {
          this.topbarText = '没有操作权限'
        } else {
          this.$message.error('查询失败')
        }
      })
    },
    clearCode () {
      this.code = ''
      this.topbarText = '新代码'
      this.topbarID = ''
      this.topbarTime = ''
    },
    submitCode () {
      if (!this.lang) {
        this.$message.warning('请选择提交语言')
        return
      }
      if (!this.code) {
        this.$message.warning('请输入代码')
        return
      }
      this.loading = true
      const params = this.$qs.stringify({
        code: this.code,
        lang: this.lang
      })
      console.log(this.lang)
      this.$axios.post(
        '/contest/' + this.cid + '/submit',
        params
      ).then(res => {
        this.loading = false
        this.$message.success('提交成功')
        this.topbarTime = this.$functions.dateTimeString(res.data.submission.created_at)
        this.topbarID = res.data.submission.id
        const statcolor = this.getStatColor(res.data.submission.status)
        this.topbarInfo.status = statcolor.stat
        this.topbarColor.status = statcolor.color
        this.topbarColor.msg = 'gray'
        this.topbarInfo.msg = res.data.submission.msg
        this.lang = res.data.submission.lang
        this.getHistory()
        // this.compareTime()
      }).catch(err => {
        this.loading = false
        this.$message.error('提交失败')
        console.log(err)
      })
    },
    showCode (sid) {
      // this.compareTime()
      this.codeLoading = true
      this.$axios.get(
        '/contest/' + this.cid + '/submission/' + sid
      ).then(res => {
        this.codeLoading = false
        this.code = res.data.contents
        this.topbarTime = this.$functions.dateTimeString(res.data.created_at)
        this.topbarID = sid
        const statcolor = this.getStatColor(res.data.status)
        this.topbarInfo.status = statcolor.stat
        this.topbarColor.status = statcolor.color
        this.topbarColor.msg = 'gray'
        this.topbarInfo.msg = res.data.msg
        this.lang = res.data.lang
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
