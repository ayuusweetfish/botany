<template>
  <v-container fluid>
    <v-snackbar
      top
      style="margin-top: 60px"
      v-model="message"
      color="error"
    >查询失败，请重试</v-snackbar>
    <v-loading-overlay v-model="loading"></v-loading-overlay>
    <v-row justify="center">
      <v-col :cols="12" :md="10" :lg="8">
        <v-card>
          <v-card-title class="title">
            <v-btn
              text
              exact
              color="primary"
              class="title pa-0"
              :to="`/contest/${$route.params.cid}/match`"
            >对局列表<v-icon>mdi-chevron-right</v-icon></v-btn>
            对局详情
          </v-card-title>
          <v-card-subtitle class="subtitle-1">对局编号：{{$route.params.mid}}</v-card-subtitle>
          <v-card-text>
            <div class="d-flex body-1 mb-4">对局状态：<status-display :status="status"></status-display></div>
            <div class="body-1 mb-2">参赛提交：共{{parties.length}}份</div>
            <div class="d-flex flex-wrap">
              <span
                class="pa-1 mr-2 mb-2 d-flex align-center"
                style="border: 1px solid silver; border-radius: 5px"
                v-for="(item, index) in parties"
                :key="index"
              >
                <router-link
                  v-if="$store.state.myrole===$consts.role.moderator"
                  :to="`/contest/${$route.params.cid}/submission/${item.id}`"
                  class="title mr-2"
                  style="text-decoration: none"
                >{{item.id}}</router-link>
                <div v-else class="title mr-2">{{item.id}}</div>
                <div class="mr-1">by</div>
                <user-tag :user="item.participant" size="small"></user-tag>
              </span>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import StatusDisplay from '../components/StatusDisplay.vue'
import UserTag from '../components/UserTag.vue'
export default {
  name: 'MatchDetail',
  components: {
    'status-display': StatusDisplay,
    'user-tag': UserTag
  },
  mounted () {
    this.getDetail()
  },
  data: () => ({
    message: false,
    scroll: 0,
    loading: false,
    parties: [],
    status: -1
  }),
  methods: {
    getDetail () {
      this.loading = true
      this.$axios.get(
        '/contest/' + this.$route.params.cid + '/match/' + this.$route.params.mid
      ).then(res => {
        this.parties = res.data.parties
        this.status = res.data.status
        this.loading = false
      }).catch(() => {
        this.loading = false
        this.message = true
      })
    }
  }
}
</script>

<style>

</style>
