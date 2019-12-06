<template>
  <div>
    <el-row style="margin-bottom: 10px">
      <el-col :span="24">
        <el-card
          class="game-title"
          :style="{backgroundImage: 'url('+bannerUrl+')'}"
          body-style="display: flex; justify-content: space-between; height: 200px; align-items: flex-end"
        >
          <div style="display: inline-flex">{{title}}</div>
          <el-button
            v-if="myRole===$consts.role.notIn && isRegOpen"
            type="primary" size="large"
            style="display: inline-flex;"
          >
            报名参加
          </el-button>
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
      this.$store.commit('enterSubSite', this.title)
      loading.close()
    }).catch(err => {
      console.log(err)
      loading.close()
      this.$message.error('查询失败，请刷新')
    })
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
  }
}
</script>

<style scoped>
  .game-title{
    font-size: 36px;
    background-size: 100% auto;
    background-position: center;
    border-radius: 5px;
  }
</style>
