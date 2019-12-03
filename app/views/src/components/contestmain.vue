<template>
  <div>
    <el-row style="margin-bottom: 10px">
      <el-col :span="24">
        <el-card class="game-title" body-style="display: flex; justify-content: space-between; height: 200px; align-items: flex-end">
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
  data () {
    return {
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
    background-image: url('../assets/demo2.png');
    background-repeat: no-repeat;
    border-radius: 5px;
  }
</style>
