<template>
  <div>
    <el-row style="margin-bottom: 10px">
      <el-col :span="24">
        <el-card
          class="contest-title"
          :style="{backgroundImage: 'url('+bannerUrl+')'}"
          body-style="display: flex; justify-content: space-between; height: 200px; align-items: flex-end"
        >
          <div style="display: inline-flex">{{title}}</div>
          <div>
            <el-button
              v-if="myRole===$consts.role.notIn && isRegOpen"
              type="primary" size="large"
              style="display: inline-flex;"
              @click="regIn"
            >
              报名参加
            </el-button>
            <el-button
              v-if="myRole===$consts.role.moderator"
              type="primary" size="large"
              style="display: inline-flex;"
              @click="goEdit"
            >
              编辑比赛
            </el-button>
            <el-button
              v-if="$store.state.privilege===$consts.privilege.superuser && !isVisible"
              type="primary" size="large"
              style="display: inline-flex;"
              @click="publishContest"
            >
              公开比赛
            </el-button>
            <el-button
              v-if="$store.state.privilege===$consts.privilege.superuser && isVisible"
              type="primary" size="large"
              style="display: inline-flex;"
              @click="hideContest"
            >
              隐藏比赛
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card>
          <el-timeline align="left">
            <el-timeline-item
              v-for="(activity, index) in events"
              :key="index"
              :timestamp="activity.time"
              :color="activity.color"
              placement="top"
            >
            <div align="left">{{activity.event}}</div>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
      <el-col :span="18">
        <el-card>
          <div align="left" style="font-size: 24px; font-weight: 600">赛事简介</div>
          <div align="left">
            <p>赛事主办方：{{owner.nickname}}</p>
            <p>{{desc}}</p>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>

export default {
  name: 'contestmain',
  created () {
    this.cid = this.$route.query.id
    this.bannerUrl = 'https://www.csgowallpapers.com/assets/images/original/mossawi_518842827528_20150625022423_816788567695.png'
    this.getContestInfo()
  },
  data () {
    return {
      cid: '',
      title: '',
      bannerUrl: '',
      isRegOpen: false,
      isVisible: false,
      desc: false,
      myRole: this.$consts.role.notIn,
      owner: '',
      events: []
    }
  },
  methods: {
    regIn () {
      const loading = this.$loading({lock: true, text: '处理中'})
      this.$axios.post(
        '/contest/' + this.cid + '/join'
      ).then(res => {
        loading.close()
        this.$message.success('报名成功')
        window.location.reload()
      }).catch(err => {
        loading.close()
        if (err.response.status !== 401) {
          this.$message('报名失败，请重试')
        }
      })
    },
    goEdit () {
      this.$router.push({
        path: '/contest_edit',
        query: {
          id: this.cid
        }
      })
    },
    publishContest () {
      const loading = this.$loading({lock: true, text: '处理中'})
      let param = this.$qs.stringify({set: true})
      this.$axios.post(
        '/contest/' + this.cid + '/publish',
        param
      ).then(res => {
        loading.close()
        this.$message.success('发布成功')
        this.getContestInfo()
      }).catch(err => {
        console.log(err)
        loading.close()
        this.$message.error('发布失败，请重试')
      })
    },
    hideContest () {
      const loading = this.$loading({lock: true, text: '处理中'})
      let param = this.$qs.stringify({set: false})
      this.$axios.post(
        '/contest/' + this.cid + '/publish',
        param
      ).then(res => {
        loading.close()
        this.$message.success('隐藏成功')
        this.getContestInfo()
      }).catch(err => {
        console.log(err)
        loading.close()
        this.$message.error('隐藏失败，请重试')
      })
    },
    getContestInfo () {
      const loading = this.$loading({lock: true, text: '查询比赛信息'})
      this.$axios.get(
        '/contest/' + this.cid + '/info'
      ).then(res => {
        this.events = []
        this.events.push({
          time: this.$functions.dateTimeString(res.data.end_time),
          event: '比赛结束',
          color: 'gray'
        })
        this.events.push({
          time: this.$functions.dateTimeString(res.data.start_time),
          event: '比赛开始',
          color: 'green'
        })
        // this.bannerUrl = res.data.banner
        this.desc = res.data.desc
        this.isRegOpen = res.data.is_reg_open
        this.isVisible = res.data.is_visible
        this.myRole = res.data.my_role
        this.owner = res.data.owner
        this.title = res.data.title
        this.$store.commit('enterSubSite', res.data)
        console.log(res.data)
        loading.close()
      }).catch(err => {
        console.log(err)
        loading.close()
        this.$message.error('查询失败，请刷新')
      })
    }
  }
}
</script>

<style scoped>
  .contest-title{
    font-size: 36px;
    background-size: 100% auto;
    background-position: center;
    border-radius: 5px;
  }
</style>
