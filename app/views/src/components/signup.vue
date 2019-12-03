<template>
  <el-container class="login-container">
    <el-header class="login-header">
      <div style="margin:auto">注册</div>
    </el-header>
    <el-main class="login-main">
      <el-form
        ref="usernamelogin"
        :model="loginInfo"
        label-suffix="left"
        label-width="0px"
        :rules="rules"
      >
        <el-row>
          <el-col :span="5">
            <div align="right" class="login-title">用户名：</div>
          </el-col>
          <el-col :span="19">
            <el-form-item prop="username" :error="loginErrUsrnm">
              <el-input
                type="text"
                v-model="loginInfo.username"
                placeholder="请输入账号"
                auto-complete="off"
                prefix-icon="el-icon-user-solid"
              ></el-input>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row>
          <el-col :span="5">
            <div align="right" class="login-title">密码：</div>
          </el-col>
          <el-col :span="19">
            <el-form-item prop="enigma" :error="loginErrPswd">
              <el-input
                type="password"
                v-model="loginInfo.enigma"
                placeholder="请输入密码"
                auto-complete="off"
                prefix-icon="el-icon-lock"
              ></el-input>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row>
          <el-col :span="5">
            <div align="right" class="login-title">邮箱：</div>
          </el-col>
          <el-col :span="19">
            <el-form-item prop="email" :error="loginErrEml">
              <el-input
                type="text"
                v-model="loginInfo.email"
                placeholder="请输入邮箱"
                auto-complete="off"
                prefix-icon="el-icon-message"
              ></el-input>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row>
          <el-col :span="5">
            <div align="right" class="login-title">验证码：</div>
          </el-col>
          <el-col :span="12">
            <el-form-item prop="enigma2" :error="loginErrPswd2">
            <el-input
              type="text"
              v-model="loginInfo.enigma2"
              placeholder="请输入验证码"
              auto-complete="off"
              prefix-icon="el-icon-s-claim"
            ></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="7">
            <img src = "../assets/demo1.png" style="width:100%; margin-top: -10px">
          </el-col>
        </el-row>
      </el-form>

      <el-row>
        <el-col :span="12">
          <el-button type="primary" @click="register" style="width: 80%">注册</el-button>
        </el-col>
        <el-col :span="12">
          <el-button @click="goLogin" style="width: 80%">返回登录</el-button>
        </el-col>
      </el-row>
    </el-main>
  </el-container>
</template>

<script>
export default {
  name: 'signin',
  created () {
    this.getCaptcha()
  },
  data () {
    // eslint-disable-next-line camelcase
    let password2_validator = (rule, value, callback) => {
      console.log(this)
      if (!value) {
        callback(new Error('请确认密码'))
      } else if (value !== this.regisInfo.password) {
        callback(new Error('密码不一致'))
      } else {
        callback()
      }
    }
    // eslint-disable-next-line camelcase
    let email_validator = (rule, value, callback) => {
      if (!value) {
        callback(new Error('请输入邮箱'))
      } else if (!/^([a-zA-Z0-9]+[-_.]?)+@([a-zA-Z0-9]+\.)+[a-z]+$/.test(value)) {
        callback(new Error('请输入格式正确的邮箱'))
      } else {
        callback()
      }
    }
    return {
      loginInfo: {
        username: '',
        enigma: '',
        enigma2: '',
        email: ''
      },
      captcha64: '',
      captchaKey: '',
      regisErrUsrnm: '',
      regisErrPswd: '',
      regisErrPswd2: '',
      regisErrEml: '',
      regisErrCpch: '',
      rules: {
        username: [
          {required: true, message: '请输入用户名', trigger: 'blur'},
          {min: 3, max: 30, message: '用户名应在3-30个字符之间', trigger: 'blur'}
        ],
        password: [
          {required: true, message: '请输入密码', trigger: 'blur'},
          {min: 5, max: 30, message: '密码应在5-30个字符之间', trigger: 'blur'}
        ],
        password2: [
          {validator: password2_validator, trigger: 'blur'}
        ],
        email: [
          {validator: email_validator, trigger: 'blur'}
        ],
        captcha: [
          {required: true, message: '请输入验证码', trigger: 'blur'},
          {min: 4, max: 4, message: '请输入4个字符', trigger: 'blur'}
        ]
      }
    }
  },
  methods: {
    getCaptcha () {
      this.$axios.get(
        '/captcha'
      ).then(res => {
        this.captcha64 = res.data.img
        this.captchaKey = res.data.kay
      // eslint-disable-next-line handle-callback-err
      }).catch(err => {
        this.$message.error('无法获取验证码，请检查网络')
      })
    },
    register () {
      this.$refs['regisform'].validate(valid => {
        if (valid) {
          this.regisErrUsrnm = ''
          this.regisErrPswd = ''
          this.regisErrPswd2 = ''
          this.regisErrEml = ''
          this.regisErrCpch = ''
          const loading = this.$loading({lock: true, text: '注册中'})
          let params = new URLSearchParams()
          params.append('username', this.regisInfo.username)
          params.append('password', this.regisInfo.password)
          params.append('email', this.regisInfo.email)
          params.append('captcha', this.regisInfo.captcha)
          this.$axios.post(
            '/register',
            params
          ).then(res => {
            loading.close()
            this.$alert('注册成功，请登录', '成功', {
              confirmButtonText: '确定',
              callback: action => {
                this.$router.replace('/')
              }
            })
          }).catch(err => {
            loading.close()
            this.$message.error('注册失败')
            if (err.response.data.error === 'wrong captcha') {
              this.regisErrCpch = '验证码错误'
              this.regisInfo.captcha = ''
              this.getCaptcha()
            } else if (err.response.data.error === 'username already exists') {
              this.regisErrUsrnm = '用户名已被注册'
              this.getCaptcha()
            } else if (err.response.data.error === 'email already exists') {
              this.regisErrEml = '邮箱已被注册'
              this.getCaptcha()
            }
          })
        }
      })
    },
    goLogin () {
      this.$router.replace('/')
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
