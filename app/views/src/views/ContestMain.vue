<template>
  <div>
    <v-loading-overlay absolute :value="loading"></v-loading-overlay>
    <v-snackbar top style="margin-top: 60px"  v-model="showErr" color="error" :timeout="3000">
      {{errMessage}}
    </v-snackbar>
    <v-img :src="bannerUrl" max-height="300" contain transition="fade-transition"></v-img>
    <div class="display-2 mb-2 d-flex justify-center">{{title}}</div>
    <v-container fluid>
      <v-row justify="center">
        <v-btn color="primary" class="ml-1 mr-1">报名参加</v-btn>
        <v-btn color="primary" class="ml-1 mr-1">公开赛事</v-btn>
        <v-btn color="primary" class="ml-1 mr-1">隐藏赛事</v-btn>
        <v-btn text color="primary" class="ml-1 mr-1" :to="`/contest/${$route.params.cid}/ranklist`">玩家排行</v-btn>
        <v-btn text color="primary" class="ml-1 mr-1" :to="`/contest/${$route.params.cid}/match`">对局列表</v-btn>
        <v-btn text color="primary" class="ml-1 mr-1" :to="`/contest/${$route.params.cid}/submission`">提交列表</v-btn>
        <v-btn text color="primary" class="ml-1 mr-1">我的代码</v-btn>
        <v-btn text color="primary" class="ml-1 mr-1" :to="`/contest/${$route.params.cid}/edit`">赛事设置</v-btn>
      </v-row>
    </v-container>
    <v-container fluid>
      <v-row justify="center">
        <v-col :cols="12" :md="9">
          <div class="body-1 mb-2"><v-icon class="mr-2 mb-1">mdi-creation</v-icon>{{owner.nickname}}</div>
          <div class="body-1 mb-4"><v-icon class="mr-2 mb-1">mdi-calendar</v-icon>{{time}}</div>
          <div class="body-1 mb-1"><v-icon>mdi-comment-processing-outline</v-icon></div>
          <div class="body-1 mb-2">{{desc}}</div>
          <v-divider></v-divider>
          <div class="body-1 mt-2" style="white-space: pre-wrap">{{details}}</div>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
export default {
  name: 'ContestMain',
  mounted () {
    this.reload()
  },
  watch: {
    '$route.params.cid': function (newval, oldval) {
      if (newval !== oldval) {
        this.reload()
      }
    }
  },
  data: () => ({
    cid: '',
    title: '',
    details: '',
    bannerUrl: '',
    isRegOpen: false,
    isVisible: false,
    desc: '',
    myRole: -1,
    owner: '',
    time: '',
    bannerLoading: false,
    loading: false,
    showErr: false,
    errMessage: '查询失败'
  }),
  methods: {
    reload () {
      this.cid = this.$route.params.cid
      this.bannerUrl = this.$axios.defaults.baseURL + '/contest/' + this.cid + '/banner'
      this.getContestInfo()
    },
    getContestInfo () {
      this.loading = true
      this.$axios.get(
        '/contest/' + this.cid + '/info'
      ).then(res => {
        const start = this.$functions.dateTimeString(res.data.end_time)
        const end = this.$functions.dateTimeString(res.data.start_time)
        this.time = start + ' TO ' + end
        this.desc = res.data.desc
        this.isRegOpen = res.data.is_reg_open
        this.isVisible = res.data.is_visible
        this.myRole = res.data.my_role
        this.owner = res.data.owner
        this.title = res.data.title
        this.details = res.data.details
        this.$store.commit('setContest', res.data)
        // this.$store.commit('enterSubSite', res.data)
        // console.log(this.myRole)
        this.loading = false
      }).catch(err => {
        if (err.response.status === 404) {
          this.errMessage = '赛事信息查询失败 (404 Not Found)'
        }
        this.loading = false
        this.showErr = true
      })
    }
  }
}
</script>

<style>

</style>
