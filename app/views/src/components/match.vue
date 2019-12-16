<template>
  <div>
    <el-row style="margin-top: 20px; margin-bottom: 20px">
      <el-card>
        <el-row>
          <el-col :span="12">
            <div class="item-title">对局人数:</div>
            <div class="item-text">{{parties.length}}</div>
          </el-col>
          <el-col :span="12">
            <div class="item-title">对局状态:</div>
            <div v-if="status === $consts.codeStat.pending" class="item-text" style="color: gray">等待中</div>
            <div v-else-if="status === $consts.codeStat.compiling" class="item-text" style="color: orange">进行中</div>
            <div v-else-if="status === $consts.codeStat.accepted" class="item-text" style="color: green">已结束</div>
            <div v-else class="item-text" style="color: red">Fatal Error</div>
          </el-col>
        </el-row>
      </el-card>
    </el-row>
    <el-row style="margin-top: 20px; margin-bottom: 20px" :gutter="10">
      <el-col v-for="(item, index) in parties" :key="index" :span="8">
        <el-card>
          <el-avatar :size="80"></el-avatar>
          <div style="font-size: 16px; font-weight: 600">{{item.participant.nickname}}</div>
          <router-link
            style="font-size: 16px; color: #409EFF; text-decoration: none"
            :to="{path: '/profile', query: {handle: item.participant.handle}}"
          >{{item.participant.handle}}</router-link>
          <div v-if="myRole===$consts.role.moderator">
            <div style="font-size: 14px; display: inline">提交ID: {{item.id}} |</div>
            <router-link
            style="font-size: 16px; color: #409EFF; text-decoration: none"
            :to="{path: '/submission_info', query: {cid: cid, sid: item.id}}"
          >查看</router-link>
          </div>
          <div v-else style="font-size: 14px">代码ID: {{item.id}}</div>
        </el-card>
      </el-col>
    </el-row>
    <el-row>
      <el-card body-style="height: 480px">
        <div align="left" style="font-size: 24px; font-weight: 600">比赛回放</div>
      </el-card>
    </el-row>
    <el-row>
      <el-card>
        <el-row>
          <el-col :span="6">
            <el-button style="width: 80%">上一回合</el-button>
          </el-col>
          <el-col :span="6">
            <el-button style="width: 80%" type="primary">暂停/开始</el-button>
          </el-col>
          <el-col :span="6">
            <el-button style="width: 80%">下一回合</el-button>
          </el-col>
          <el-col :span="6">
            <el-button style="width: 80%" type="text">下载录像文件</el-button>
          </el-col>
        </el-row>
      </el-card>
    </el-row>
  </div>
</template>

<script>
export default {
  name: 'match',
  created () {
    this.cid = this.$route.query.cid
    this.mid = this.$route.query.mid
    this.getInfo()
  },
  data () {
    return {
      parties: [],
      mid: '',
      cid: '',
      status: 0,
      myRole: -1,
    }
  },
  methods: {
    getInfo() {
      this.parties = []
      const loading = this.$loading({lock: true, text: '加载中'})
      this.$axios.get(
        '/contest/' + this.cid + '/match/' + this.mid
      ).then(res => {
        console.log(res.data)
        this.parties = res.data.parties
        this.status = res.data.status
        this.$axios.get(
          '/contest/' + this.cid + '/info'
        ).then(info => {
          loading.close()
          this.myRole = info.data.my_role
        })
      }).catch(err => {
        console.log(err)
        loading.close()
        this.$message.error('加载失败')
      })
    }
  }
}
</script>

<style scoped>
.item-title{
  display: inline;
  color: gray;
  font-weight: 600;
}
.item-text{
  display: inline;
  color: black;
  font-weight: 400;
}
</style>