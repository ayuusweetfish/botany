<template>
  <v-container fluid>
    <v-snackbar
      top
      style="margin-top: 60px"
      v-model="message"
      color="error"
    >查询失败，请重试</v-snackbar>
    <v-loading-overlay v-model="loading"></v-loading-overlay>
    <v-fab-transition>
      <v-btn
        fixed
        bottom
        right
        fab
        v-scroll="onScroll"
        v-show="scroll>100"
        color="primary"
        class="mr-6 mb-6"
        @click="$vuetify.goTo(0)"
      ><v-icon>mdi-chevron-up</v-icon></v-btn>
    </v-fab-transition>
    <v-fab-transition>
      <v-btn
        fixed
        bottom
        right
        fab
        v-scroll="onScroll"
        v-show="scroll===0"
        color="primary"
        class="mr-6 mb-6"
        @click="$vuetify.goTo($refs.code)"
      ><v-icon>mdi-chevron-down</v-icon></v-btn>
    </v-fab-transition>
    <v-row justify="center">
      <v-col :cols="12" :md="10" :lg="8">
        <v-card>
          <v-card-title class="title">
            <v-btn
              text
              exact
              color="primary"
              class="title pa-0"
              :to="`/contest/${$route.params.cid}/submission`"
            >提交列表<v-icon>mdi-chevron-right</v-icon></v-btn>
            提交详情
          </v-card-title>
          <v-card-subtitle class="subtitle-1">第{{$route.params.sid}}号提交@{{$store.state.cname}}</v-card-subtitle>
          <v-card-text>
            <div>
              <user-tag :user="participant"></user-tag>
            </div>
            <div class="body-1">代码语言：{{lang}}</div>
            <div class="d-flex body-1 mb-4">
              代码状态：<status-display mode="submission" :status="status"></status-display>
            </div>
            <div>
              <code-edit
                :value="info"
                height="auto"
                lang="text"
                theme="material"
                :options="cmOptions"
              ></code-edit>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col :cols="12" :md="10" :lg="8">
        <v-card outlined ref="code">
          <code-edit
            :value="code"
            height="auto"
            :lang="getHighlightLang()"
            :options="cmOptions"
          ></code-edit>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import CodeEditor from '../components/CodeEditor.vue'
import StatusDisplay from '../components/StatusDisplay.vue'
import UserTag from '../components/UserTag.vue'
export default {
  name: 'SubmissionDetail',
  components: {
    'code-edit': CodeEditor,
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
    participant: {},
    code: '',
    lang: '',
    status: -1,
    info: '',
    cmOptions: {
      lineNumbers: true,
      indentUnit: 2,
      tabSize: 2,
      autoCloseBrackets: true,
      readOnly: true
    }
  }),
  methods: {
    onScroll () {
      this.scroll = window.pageYOffset
    },
    getDetail () {
      this.loading = true
      this.$axios.get(
        '/contest/' + this.$route.params.cid + '/submission/' + this.$route.params.sid
      ).then(res => {
        this.participant = res.data.participant
        this.code = res.data.contents.slice(0, -1)
        this.lang = res.data.lang
        this.info = res.data.msg.slice(0, -1)
        this.status = res.data.status
        this.loading = false
      }).catch(() => {
        this.loading = false
        this.message = true
      })
    },
    getHighlightLang () {
      if (this.lang.endsWith('.c') || this.lang.endsWith('cpp')) {
        return 'c'
      } else if (this.lang.endsWith('lua')) {
        return 'lua'
      } else if (this.lang.endsWith('py')) {
        return 'python'
      }
      return 'text'
    }
  }
}
</script>

<style>

</style>
