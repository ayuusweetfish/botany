<template>
  <v-card
    :raised="$vuetify.breakpoint.mdAndUp"
    :outlined="!$vuetify.breakpoint.mdAndUp"
    :width="$vuetify.breakpoint.mdAndUp? '80%' : '100%'"
    style="border: none"
  >
    <v-dialog
      v-model="showDialog"
      max-width="480"
    >
      <v-card style="border-left: 10px solid green">
        <v-card-title class="headline success--text">注册成功！</v-card-title>
        <v-card-text>
          <div class="title">欢迎加入BotAny，{{nickname}}！</div>
          <div class="title">你的账户名为：{{handle}}</div>
          <div class="title">是否登录这个账号？</div>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="primary" class="title ma-4" @click="login" :loading="logging" :disabled="logging">登录</v-btn>
          <v-btn text color="secondary" class="title ma-4" @click="showDialog=false">取消</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-snackbar top style="margin-top: 60px" v-model="showErr" color="error" :timeout="3000">
      {{errMessage}}<v-btn text @click="showErr=false">确定</v-btn>
    </v-snackbar>
    <v-card-title><div class="headline pl-6 pr-6 pb-4 pt-4">加入BotAny</div></v-card-title>
    <v-card-text>
      <v-form
        ref="signupForm"
        v-model="valid"
      >
        <v-text-field
          v-model="handle"
          label="账号名称"
          :rules="[validators.required, validators.minmax, v=>!errHandle||errHandle]"
          class="pl-6 pr-6"
          hint="3-30个字符，用于登录"
          counter
        ></v-text-field>
        <v-text-field
          v-model="password"
          label="密码"
          :rules="[validators.required, validators.minmax, considerP2]"
          :append-icon="showPassword? 'mdi-eye' : 'mdi-eye-off'"
          @click:append="showPassword=!showPassword"
          :type="showPassword? 'text' : 'password'"
          class="pl-6 pr-6"
          hint="3-30个字符"
          counter
        ></v-text-field>
        <v-text-field
          ref="password2"
          v-model="password2"
          label="确认密码"
          :rules="[validators.required, v=>v===password||'密码输入不一致']"
          :append-icon="showPassword2? 'mdi-eye' : 'mdi-eye-off'"
          @click:append="showPassword2=!showPassword2"
          :type="showPassword? 'text' : 'password'"
          class="pl-6 pr-6"
          hint="重复上面的密码"
          counter
        ></v-text-field>
        <v-text-field
          v-model="email"
          label="邮箱"
          :rules="[validators.required, validators.email, v=>!errEmail||errEmail]"
          class="pl-6 pr-6"
          hint="请输入一个与账号绑定的邮箱"
          counter
        ></v-text-field>
        <v-text-field
          v-model="nickname"
          label="昵称"
          :rules="[validators.required, validators.minmax]"
          class="pl-6 pr-6"
          hint="输入你在botany显示的名称"
          counter
        ></v-text-field>
        <v-text-field
          v-model="captcha"
          label="验证码"
          :rules="[validators.required, validators.captcha, v=>!errCaptcha||errCaptcha]"
          class="pl-6 pr-6 mt-6"
          :counter="4"
          :loading="captchaLoading"
          outlined
        >
          <template v-slot:append>
            <v-img :src="captcha64" @click="getCaptcha" style="cursor: pointer"></v-img>
          </template>
        </v-text-field>
      </v-form>
    </v-card-text>
    <v-card-actions>
      <v-row class="pl-4 pr-4">
        <v-col :cols="12" :md="6">
          <v-btn
            color="primary"
            block large
            :disabled="!valid||loading"
            @click="signup"
            :loading="loading"
          >注册账号
          </v-btn>
        </v-col>
        <v-col :cols="12" :md="6">
          <v-btn text
            color="secondary"
            block large
            :to="{path: '/register/login', query: {redirect: $route.query.redirect}}"
          >返回登录
        </v-btn>
        </v-col>
      </v-row>
    </v-card-actions>
  </v-card>
</template>

<script>
export default {
  name: 'signup',
  mounted () {
    this.getCaptcha()
  },
  data: el => ({
    showErr: false,
    errMessage: '',
    showDialog: false,
    logging: false,
    validators: {
      required: v => !!v || '输入不能为空',
      minmax: v => (v.length >= 3 && v.length <= 30) || '长度不符合要求',
      email: v => (/^([a-zA-Z0-9]+[-_.]?)+@([a-zA-Z0-9]+\.)+[a-z]+$/.test(v)) || '请输入格式合法的邮箱',
      captcha: v => v.length === 4 || '验证码格式错误'
    },
    considerP2: v => {
      if (el.password2) {
        el.$refs.password2.validate()
      }
      return true
    },
    showPassword: false,
    showPassword2: false,
    handle: '',
    password: '',
    password2: '',
    email: '',
    nickname: '',
    captcha: '',
    captcha64: '',
    captchaKey: '',
    errHandle: '',
    errEmail: '',
    errCaptcha: '',
    captchaLoading: false,
    valid: false,
    loading: false
  }),
  methods: {
    getCaptcha () {
      this.captchaLoading = true
      this.$axios.get(
        '/captcha'
      ).then(res => {
        this.captcha64 = res.data.img
        this.captchaKey = res.data.key
        this.captchaLoading = false
      }).catch(() => {
        this.captchaLoading = false
      })
    },
    signup () {
      if (!this.$refs.signupForm.validate()) {
        return
      }
      this.loading = true
      const params = this.$qs.stringify({
        handle: this.handle,
        password: this.password,
        email: this.email,
        nickname: this.nickname,
        captcha_value: this.captcha,
        captcha_key: this.captchaKey
      })
      this.$axios.post(
        '/signup',
        params
      ).then(res => {
        this.loading = false
        this.showDialog = true
      }).catch(err => {
        this.loading = false
        this.getCaptcha()
        err.response.data.err.forEach(item => {
          switch (parseInt(item)) {
            case 1:
              this.errHandle = '用户名已被注册'
              break
            case 2:
              this.errEmail = '邮箱已被注册'
              break
            case 3:
              this.errCaptcha = '验证码不正确'
              break
            default:
              break
          }
        })
        this.$refs.signupForm.validate()
        this.errHandle = ''
        this.errEmail = ''
        this.errCaptcha = ''
        this.errMessage = '注册失败，请检查表单'
        this.showErr = true
      })
    },
    login () {
      this.logging = true
      const params = this.$qs.stringify({
        handle: this.handle,
        password: this.password
      })
      this.$axios.post(
        '/login',
        params
      ).then(res => {
        const logindata = {
          id: res.data.id,
          handle: res.data.handle,
          privilege: parseInt(res.data.privilege),
          nickname: res.data.nickname
        }
        this.$store.commit('login', logindata)
        this.logging = false
        if (this.$route.query.redirect) {
          this.$router.push(this.$store.state.redirect)
        } else {
          this.$router.push('/')
        }
      }).catch(() => {
        this.logging = false
        this.errMessage = '登录失败，请重试'
      })
    }
  }
}
</script>

<style>

</style>
