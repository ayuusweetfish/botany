<template>
  <div>
    <v-snackbar
      style="margin-top: 60px"
      top
      v-model="showMessage"
      :color="messageType"
      :timeout="messageType==='success'? 0 : 3000"
    >
      {{message}}
      <v-btn outlined color="white" v-if="messageType==='success'" @click="goNext">确定</v-btn>
    </v-snackbar>
    <v-loading-overlay v-model="loading"></v-loading-overlay>
    <v-speed-dial
      fixed
      v-model="fab"
      bottom right
      direction="top"
      transition="slide-y-reverse-transition"
      class="mb-11"
    >
      <template v-slot:activator>
        <v-btn
          v-model="fab"
          fab
          color="primary darken-3"
          dark
        >
          <v-icon v-if="fab">mdi-close</v-icon>
          <v-icon v-else>mdi-format-list-checkbox</v-icon>
        </v-btn>
      </template>
      <v-btn small fab color="primary" @click="$vuetify.goTo($refs.end)">
        <v-icon>mdi-chevron-down</v-icon>
      </v-btn>
      <v-btn small fab color="primary" @click="$vuetify.goTo($refs.script)">
        <v-icon>mdi-script-text-outline</v-icon>
      </v-btn>
      <v-btn small fab color="primary" @click="$vuetify.goTo($refs.form)">
        <v-icon>mdi-chevron-up</v-icon>
      </v-btn>
    </v-speed-dial>
    <v-form ref="form">
      <v-row>
        <v-col :cols="12" :md="6">
          <v-text-field
            label="名称"
            v-model="title"
            @change="setStall"
            prepend-icon="mdi-tag-outline"
            :counter="20"
            :rules="rules.title"
          ></v-text-field>
        </v-col>
        <v-col :cols="12" :md="6">
          <v-radio-group row label="开放报名" v-model="isRegOpen" @change="setStall" prepend-icon="mdi-key-outline">
            <v-radio label="是" :value="true"></v-radio>
            <v-radio label="否" :value="false"></v-radio>
          </v-radio-group>
        </v-col>
      </v-row>
      <v-row>
        <v-col :cols="12" :md="6">
          <date-time-picker
            v-model="dateStart"
            prepend-icon="mdi-calendar-today"
            label="开始日期"
            required
            @change="setStall"
            ref="dateStart"
          ></date-time-picker>
        </v-col>
        <v-col :cols="12" :md="6">
          <date-time-picker
            v-model="dateEnd"
            prepend-icon="mdi-calendar"
            label="结束日期"
            @change="setStall"
            required
            :min="dateStart"
            ref="dateEnd"
          ></date-time-picker>
        </v-col>
      </v-row>
      <user-select v-model="moderators" @change="setStall"></user-select>
      <v-textarea
        label="简介"
        outlined
        v-model="desc"
        @change="setStall"
        prepend-icon="mdi-comment-processing-outline"
        :counter="255"
        :rows="3"
        :rules="rules.desc"
      >
      </v-textarea>
      <v-textarea
        label="详细介绍"
        v-model="details"
        @change="setStall"
        :counter="4095"
        outlined
        :rows="12"
        prepend-icon="mdi-comment-text-outline"
        :rules="rules.details"
      >
      </v-textarea>
      <div ref="script" class="d-flex align-start mb-8">
        <v-icon class="mr-2">mdi-script-text-outline</v-icon>
        <div style="width: 100%">
          <div class="body-1 secondary--text mb-2">赛制脚本</div>
          <v-card outlined>
            <code-editor
              v-model="script"
              lang="lua"
              height="480px"
              theme="material"
              :options="cmOptions"
              @change="setStall"
            ></code-editor>
          </v-card>
        </div>
      </div>
      <div class="d-flex align-start mb-8">
        <v-icon class="mr-2">mdi-iframe-outline</v-icon>
        <div style="width: 100%">
          <div class="body-1 secondary--text mb-2">回放模板</div>
          <v-card outlined>
            <code-editor
              v-model="playback"
              lang="html"
              height="480px"
              theme="material"
              :options="cmOptions"
              @change="setStall"
            ></code-editor>
          </v-card>
        </div>
      </div>
    </v-form>
    <v-row justify="space-around" ref="end">
      <v-col :cols="4">
        <v-btn block color="primary" @click="confirm" :disabled="submiting" :loading="submiting">提交</v-btn>
      </v-col>
      <v-col :cols="4">
        <v-btn block outlined @click="cancel">取消</v-btn>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import UserSelector from '../components/UserSelector.vue'
import DateTimePicker from '../components/DateTimePicker.vue'
import CodeEditor from '../components/CodeEditor.vue'
export default {
  props: {
    mode: {
      type: String,
      default: 'edit'
    },
    cid: String
  },
  components: {
    'date-time-picker': DateTimePicker,
    'user-select': UserSelector,
    'code-editor': CodeEditor
  },
  mounted () {
    this.$store.commit('setStall', false)
    if (this.mode === 'edit') {
      this.loadFlag = false
      this.htmlLoadFlag = false
      this.loadInfo()
    } else {
      this.loadFlag = true
      this.htmlLoadFlag = true
    }
  },
  data: () => ({
    valid: false,
    loading: false,
    loadFlag: false,
    htmlLoadFlag: false,
    showMessage: false,
    message: '',
    messageType: '',
    submiting: false,
    title: '',
    dateStart: '',
    dateEnd: '',
    desc: '',
    details: '',
    isRegOpen: false,
    moderators: [],
    script: '',
    playback: '',
    cmOptions: {
      lineWrapping: true,
      tabSize: 2,
      lineNumbers: true,
      line: true
    },
    fab: false,
    rules: {
      title: [
        v => !!v || '请输入标题',
        v => v.length <= 20 || '标题字数超过限制'
      ],
      dateStart: [
        v => !!v || '请选择开始日期'
      ],
      dateEnd: [
        v => !!v || '请选择结束日期'
      ],
      desc: [
        v => !!v || '请输入赛事简介',
        v => v.length <= 255 || '字数超过限制'
      ],
      details: [
        v => !!v || '请输入详细介绍',
        v => v.length <= 4095 || '字数超过限制'
      ]
    }
  }),
  methods: {
    loadInfo () {
      this.loading = true
      this.$axios.get(
        '/contest/' + this.$route.params.cid + '/info'
      ).then(res => {
        this.$store.commit('setContest', res.data)
        if (res.data.my_role !== this.$consts.role.moderator) {
          this.message = '没有权限进行这个操作'
          this.messageType = 'error'
          this.showMessage = true
          return
        }
        this.script = res.data.script
        this.title = res.data.title
        this.dateStart = this.$functions.dateTimeString(new Date(res.data.start_time))
        this.dateEnd = this.$functions.dateTimeString(new Date(res.data.end_time))
        this.desc = res.data.desc
        this.details = res.data.details
        this.isRegOpen = res.data.is_reg_open
        // console.log(this.moderators)
        this.moderators = res.data.moderators
        this.loadFlag = true
        this.loading = false
      }).catch(() => {
        this.loading = false
        // push 404
        this.message = '查询失败'
        this.messageType = 'error'
        this.showMessage = true
      })
      this.$axios.get(
        '/contest/' + this.cid + '/match/0/playback'
      ).then(res => {
        this.playback = res.data
        this.htmlLoadFlag = true
      }).catch(() => {})
    },
    setStall () {
      if (this.loadFlag && this.htmlLoadFlag) {
        this.$store.commit('setStall', true)
      }
    },
    searchUsers () {
      if (this.searchText !== '') {
        this.searchLoading = true
        this.$axios.get(
          '/user_search/' + this.searchText
        ).then(res => {
          this.selections = res.data
          this.searchLoading = false
        }).catch(() => {
          this.searchLoading = false
        })
      } else {
        this.selections = []
      }
    },
    isSelectionDisabled (id) {
      for (let i = 0; i < this.moderators.length; ++i) {
        if (this.moderators[i].id === id) {
          return true
        }
      }
      return false
    },
    removeModerator (id) {
      this.moderators.splice(this.moderators.findIndex(item => item.id === id), 1)
    },
    confirm () {
      const v1 = this.$refs.dateStart.validate()
      const v2 = this.$refs.dateEnd.validate()
      const v3 = this.$refs.form.validate()
      if (v1 && v2 && v3) {
        this.submiting = true
        let params = {
          title: this.title,
          start_time: Math.round(new Date(this.dateStart).getTime() / 1000),
          end_time: Math.round(new Date(this.dateEnd).getTime() / 1000),
          desc: this.desc,
          details: this.details,
          is_reg_open: this.isRegOpen,
          script: this.script,
          playback: this.playback
        }
        if (this.moderators.length > 0) {
          const idList = Array.from(this.moderators, item => item.id)
          params.moderators = idList.join(',')
        }
        params = this.$qs.stringify(params)
        const tail = this.mode === 'edit' ? `${this.$route.params.cid}/edit` : '/create'
        this.$axios.post(
          '/contest/' + tail,
          params
        ).then(res => {
          this.$store.commit('setStall', false)
          this.submiting = false
          this.messageType = 'success'
          this.message = '提交成功'
          this.showMessage = true
        }).catch(() => {
          this.submiting = false
          this.messageType = 'error'
          this.message = '提交失败，请检查表单后重试'
          this.showMessage = true
        })
      } else {
        this.messageType = 'warning'
        this.message = '表单未正确填写，请检查后重试'
        this.showMessage = true
      }
    },
    goNext () {
      const url = this.mode === 'edit' ? `/contest/${this.$route.params.cid}/main` : '/'
      this.$router.push(url)
    },
    cancel () {
      const url = this.mode === 'edit' ? `/contest/${this.$route.params.cid}/main` : '/'
      this.$router.push(url)
    }
  }
}
</script>

<style>

</style>
