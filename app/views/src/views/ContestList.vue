<template>
  <div>
    <v-snackbar top style="margin-top: 60px"  v-model="showErr" color="error" :timeout="3000">
      {{errMessage}}
    </v-snackbar>
    <v-container fluid>
      <h1>欢迎来到BotAny</h1>
      <div>当前共有{{total}}场赛事</div>
      <v-row>
        <v-col
          :cols="12" :md="6"
          v-for="(item, index) in contests" :key="index"
        >
          <v-card
            :to="`/contest/${item.id}/main`"
          >
            <v-img :src="$axios.defaults.baseURL + '/contest/' + item.id + '/banner'" :height="180" contain></v-img>
            <v-card-title>
              <div class="d-inline">{{item.title}}</div>
              <div class="d-inline primary--text ml-2" v-if="item.my_role===$consts.role.moderator">我管理的比赛</div>
            </v-card-title>
            <v-card-subtitle>{{item.time}}</v-card-subtitle>
            <v-card-text>
              <div>{{item.desc}}</div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
export default {
  name: 'ContestList',
  data: () => ({
    total: 0,
    contests: [],
    errMessage: '无法连接服务器，请检查网络',
    showErr: false
  }),
  mounted () {
    this.getContestList()
  },
  methods: {
    getContestList () {
      this.$axios.get('/contest/list').then(res => {
        this.contests = res.data
        this.total = res.data.length
        this.contests.forEach(item => {
          const start = this.$functions.dateTimeString(item.start_time)
          const end = this.$functions.dateTimeString(item.end_time)
          item.time = start + ' TO ' + end
          console.log(item)
        })
      }).catch(err => {
        if (err.response.state >= 400) {
          this.showErr = true
        }
      })
    }
  }
}
</script>

<style>

</style>
