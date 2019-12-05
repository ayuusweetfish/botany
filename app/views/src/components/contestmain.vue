<template>
  <div>
    <el-row style="margin-bottom: 10px">
      <el-col :span="24">
        <el-card 
          class="game-title"
          :style="{backgroundImage: 'url('+bannerUrl+')'}"
          body-style="display: flex; justify-content: space-between; height: 200px; align-items: flex-end"
          >
          <div style="display: inline-flex">GStrategy by SSAST</div>
          <el-button type="primary" size="large" style="display: inline-flex;">报名参加</el-button>
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
            <p>赛事主办方：软件学院科协</p>
            <p>决赛赛制：双败淘汰制</p>
            <p>**more info**</p>
            <p>**more info**</p>
            <p>**more info**</p>
            <p>**more info**</p>
            <p>**more info**</p>
            <p>**more info**</p>
            <p>**more info**</p>
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
    this.bannerUrl = 'https://www.csgowallpapers.com/assets/images/original/mossawi_518842827528_20150625022423_816788567695.png'
    const loading = this.$loading({lock: true, text:'查询比赛信息'})
    this.$axios.get(
      '/contest/' + this.$route.query.id + '/info'
    ).then(res=>{
      
    })
    this.$store.commit('setSiteName', )
  },
  data () {
    return {
      cid: '',
      title: '',
      bannerUrl: '',
      events: [
        {
          time: '2019-01-01',
          event: '总决赛',
          color: 'gray'
        },
        {
          time: '2019-11-01',
          event: '排名赛开始',
          color: 'gray'
        },
        {
          time: '2019-10-01',
          event: '可以开始非排名比赛',
          color: 'green'
        },
        {
          time: '2019-09-01',
          event: '比赛资源上传',
          color: 'green'
        }
      ]
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
