<template>
  <el-container class="login-container">
    <el-header class="login-header">
      <div style="margin:auto">登录</div>
    </el-header>
    <el-main class="login-main">
      <el-form
        ref="loginform"
        :model="loginInfo"
        label-position="right"
        label-width="100px"
        hide-required-asterisk
        :rules="rules"
      >
        <el-row>
          <el-form-item prop="handle" :error="loginErrHandle" label="用户名：">
            <el-input
              type="text"
              v-model="loginInfo.handle"
              placeholder="请输入账号"
              auto-complete="off"
              prefix-icon="el-icon-user-solid"
            ></el-input>
          </el-form-item>
        </el-row>

        <el-row>
          <el-form-item prop="password" :error="loginErrPswd" label="密码：">
            <el-input
              type="password"
              v-model="loginInfo.password"
              placeholder="请输入密码"
              auto-complete="off"
              prefix-icon="el-icon-lock"
            ></el-input>
          </el-form-item>
        </el-row>

        <!-- <el-row>
          <el-col :span="5">
            <div align="right" class="login-title">验证码：</div>
          </el-col>
          <el-col :span="12">
            <el-form-item prop="captcha" :error="loginErrCpch">
            <el-input
              type="text"
              v-model="loginInfo.captcha"
              placeholder="请输入验证码"
              auto-complete="off"
              prefix-icon="el-icon-s-claim"
            ></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="7">
            <img :src="captcha64" style="width:100%; margin-top: -10px">
          </el-col>
        </el-row>-->
      </el-form>

      <el-row>
        <el-col :span="12">
          <el-button type="primary" @click="login" style="width: 80%">登录</el-button>
        </el-col>
        <el-col :span="12">
          <el-button @click="goSignup" style="width: 80%">注册</el-button>
        </el-col>
      </el-row>
      <!--<el-row>
        <el-button type="text">忘记密码？</el-button>
      </el-row>-->
    </el-main>
  </el-container>
</template>

<script>
export default {
  name: 'login',
  created () {
    if (this.$route.query && this.$route.query.redirect && this.$store.state.afterLogin) {
      this.nextRoute = this.$store.state.afterLogin
      this.redirect = this.$route.query.redirect
    }
  },
  data () {
    return {
      loginInfo: {
        handle: '',
        password: ''
        // captcha: ''
      },
      captcha64: '',
      captchaKey: '',
      loginErrHandle: '',
      loginErrPswd: '',
      loginErrCpch: '',
      nextRoute: {
        path: '/'
      },
      redirect: false,
      rules: {
        handle: [
          {required: true, message: '请输入账号', trigger: 'blur'},
          {max: 30, message: '输入过长', trigger: 'blur'}
        ],
        password: [
          {required: true, message: '请输入密码', trigger: 'blur'},
          {max: 30, message: '输入过长', trigger: 'blur'}
        ]
        // captcha: [
        //   {required: true, message: '请输入验证码', trigger: 'blur'},
        //   {min: 4, max: 4, message: '请输入4个字符', trigger: 'blur'}
        // ]
      }
    }
  },
  methods: {
    getCaptcha () {
      this.$axios.get(
        '/captcha'
      ).then(res => {
        this.captcha64 = res.data.img
        this.captchaKey = res.data.key
      // eslint-disable-next-line handle-callback-err
      }).catch(err => {
        this.$message.error('无法获取验证码，请检查网络')
      })
    },
    login () {
      this.$refs['loginform'].validate(valid => {
        if (valid) {
          this.loginErrHandle = ''
          this.loginErrPswd = ''
          this.loginErrCpch = ''
          const loading = this.$loading({lock: true, text: '登录中'})
          let params = this.$qs.stringify({
            'handle': this.loginInfo.handle,
            'password': this.loginInfo.password
          })
          this.$axios.post(
            '/login',
            params
          ).then(res => {
            let logindata = {
              'id': res.data.id,
              'handle': res.data.handle,
              'privilege': parseInt(res.data.privilege),
              'nickname': res.data.nickname
            }
            this.$store.commit('login', logindata)
            loading.close()
            console.log(this.nextRoute)
            this.$router.push(this.nextRoute)
          }).catch(err => {
            loading.close()
            console.log(err)
            this.$message.error('登录失败')
            this.loginErrUsrnm = '用户名或密码错误'
            this.loginErrPswd = '用户名或密码错误'
            // if (err.response.data.error === 'wrong captcha') {
            //   this.loginErrCpch = '验证码错误'
            //   this.loginInfo.captcha = ''
            //   this.getCaptcha()
            // } else if (err.response.data.error === 'wrong username or password') {
            //   this.loginErrUsrnm = '用户名或密码错误'
            //   this.loginErrPswd = '用户名或密码错误'
            //   this.loginInfo.password = ''
            //   this.getCaptcha()
            // }
          })
        }
      })
    },
    goSignup () {
      this.$router.push({
        path: '/signup',
        query: {
          redirect: this.redirect
        }
      })
    }
  }
}
</script>

<style scoped>
.login-container {
  border-radius: 5px;
  background-clip: padding-box;
  margin: 100px auto;
  padding: 20px;
  width: 550px;
  border: 1px solid silver;
  font-weight: 600;
  font-size: 25px;
}
.login-header {
  display: flex;
  text-align: center;
  border-top: none;
  border-left: none;
  border-right: none;
  border-bottom: 1px solid silver;
}
.login-main {
  font-weight: 400;
  font-size: 18px;
}
.login-title {
  font-weight: 400;
  font-size: 18px;
  vertical-align: middle;
  line-height: 40px;
}
.login-tip {
  font-weight: 400;
  font-size: 16px;
  color: silver;
  margin-bottom: 0;
  padding: 0;
  line-height: 22px;
}
</style>
