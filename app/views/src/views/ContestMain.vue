<template>
  <div>
    <v-loading-overlay absolute :value="loading"></v-loading-overlay>
    <v-snackbar top style="margin-top: 60px"  v-model="showErr" color="error" :timeout="3000">
      {{errMessage}}
    </v-snackbar>
    <v-snackbar top style="margin-top: 60px"  v-model="showSuccess" color="success" :timeout="3000">
      {{successMsg}}
    </v-snackbar>
    <v-img :src="bannerUrl" max-height="300" transition="fade-transition"></v-img>
    <input
      ref="banner-upload"
      type="file"
      style="display: none"
      accept="image/gif, image/jpeg, image/png"
      min="1"
      max="1"
      @change="uploadBanner"/>
    <div class="display-2 mt-6 mb-2 d-flex justify-center">{{title}}</div>
    <v-container fluid>
      <v-row justify="center">
        <v-btn color="primary"
          class="ml-1 mr-1"
          v-if="checkAuth('notIn')"
          :loading="joining"
          :disabled="joining"
          @click="join"
        >报名参加</v-btn>
        <v-btn color="primary"
          class="ml-1 mr-1"
          v-if="checkAuth('moderator')"
          :loading="bannerLoading"
          :disabled="bannerLoading"
          @click="changeBanner"
        >更改Banner</v-btn>
        <div v-if="$store.state.privilege === $consts.privilege.superuser">
          <v-btn
            color="primary" class="ml-1 mr-1"
            v-if="!isVisible"
            @click="setContest(true)"
            :loading="setting"
            :disabled="setting"
          >公开赛事</v-btn>
          <v-btn
            color="primary" class="ml-1 mr-1"
            v-else
            @click="setContest(false)"
            :loading="setting"
            :disabled="setting"
          >隐藏赛事</v-btn>
        </div>
        <v-btn
          text color="primary" class="ml-1 mr-1"
          :to="`/contest/${$route.params.cid}/script`"
          v-if="checkAuth('moderator')"
        >脚本操作</v-btn>
        <v-btn
          text color="primary" class="ml-1 mr-1"
          :to="`/contest/${$route.params.cid}/judge#submit`"
          v-if="checkAuth('moderator')"
        >设置裁判</v-btn>
        <v-btn
          text color="primary" class="ml-1 mr-1"
          :to="`/contest/${$route.params.cid}/edit`"
          v-if="checkAuth('moderator')"
        >赛事编辑</v-btn>
        <v-btn
          text color="primary" class="ml-1 mr-1"
          :to="`/contest/${$route.params.cid}/submit#edit`"
          v-if="checkAuth('imIn')"
        >我的代码</v-btn>
        <v-btn
          text color="primary" class="ml-1 mr-1"
          :to="`/contest/${$route.params.cid}/ranklist`"
        >玩家排行</v-btn>
        <v-btn
          text color="primary" class="ml-1 mr-1"
          :to="`/contest/${$route.params.cid}/match`"
        >对局列表</v-btn>
        <v-btn
          text color="primary" class="ml-1 mr-1"
          :to="`/contest/${$route.params.cid}/submission`"
        >提交列表</v-btn>
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
    start: new Date(),
    end: new Date(),
    time: '',
    bannerLoading: false,
    loading: false,
    showErr: false,
    errMessage: '查询失败',
    showSuccess: false,
    successMsg: '',
    joining: false,
    setting: false
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
        const end = this.$functions.dateTimeString(res.data.end_time)
        const start = this.$functions.dateTimeString(res.data.start_time)
        this.start = new Date(start)
        this.end = new Date(end)
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
    },
    checkAuth (type) {
      return this.$store.state.myrole === this.$consts.role[type]
    },
    join () {
      if (this.$store.state.handle === '') {
        this.$router.push({ path: '/register/login', query: { redirect: true } })
      }
      this.joining = true
      this.$axios.post(
        '/contest/' + this.cid + '/join'
      ).then(res => {
        this.successMsg = '成功加入赛事'
        this.showSuccess = true
        this.reload()
      }).catch(() => {
        this.errMessage = '加入赛事失败'
        this.showErr = true
      })
    },
    setContest (visible) {
      this.setting = true
      const param = this.$qs.stringify({ set: visible })
      this.$axios.post(
        '/contest/' + this.cid + '/publish',
        param
      ).then(res => {
        this.setting = false
        this.successMsg = '设置成功'
        this.showSuccess = true
        this.reload()
      }).catch(() => {
        this.setting = false
        this.errMessage = '设置失败'
        this.showErr = true
      })
    },
    changeBanner () {
      this.$refs['banner-upload'].click()
    },
    uploadBanner () {
      const files = this.$refs['banner-upload'].files
      if (!files || files.length !== 1) {
        this.errMessage = '上传文件过多'
        this.showErr = true
        return
      }
      if (files[0].size >= 1024 * 1024) {
        this.errMessage = '上传文件过大'
        this.showErr = true
        return
      }
      const namelist = files[0].name.split('.')
      const filetype = namelist[namelist.length - 1]
      if (['jpg', 'jpeg', 'gif', 'png'].indexOf(filetype) === -1) {
        this.errMessage = '上传文件格式错误'
        this.showErr = true
        return
      }
      this.bannerLoading = true
      const bannerpckt = new FormData()
      bannerpckt.append('file', files[0])
      this.$axios.post(
        '/contest/' + this.cid + '/banner/upload',
        bannerpckt
      ).then(res => {
        this.bannerLoading = false
        window.location.reload()
      }).catch(() => {
        this.bannerLoading = false
        this.errMessage = '上传失败'
        this.showErr = true
      })
    }
  }
}
</script>

<style>

</style>
