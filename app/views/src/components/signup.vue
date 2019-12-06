<template>
  <el-container class="login-container">
    <el-header class="login-header">
      <div style="margin:auto">注册</div>
    </el-header>
    <el-main class="login-main">
      <el-form
        ref="signupform"
        :model="signupInfo"
        label-width="100px"
        :rules="rules"
      >
        <el-row>
          <el-form-item prop="handle" :error="signupErrUsrnm" label="用户名：">
            <el-input
              type="text"
              v-model="signupInfo.handle"
              placeholder="输入登录时使用的账户名称"
              auto-complete="off"
              prefix-icon="el-icon-user-solid"
            ></el-input>
          </el-form-item>
        </el-row>

        <el-row>
          <el-form-item prop="password" :error="signupErrPswd" label="密码：">
            <el-input
              type="password"
              v-model="signupInfo.password"
              placeholder="请输入密码"
              auto-complete="off"
              prefix-icon="el-icon-lock"
            ></el-input>
          </el-form-item>
        </el-row>

        <el-row>
          <el-form-item prop="password2" :error="signupErrPswd2" label="确认密码：">
            <el-input
              type="password"
              v-model="signupInfo.password2"
              placeholder="再次输入密码"
              auto-complete="off"
              prefix-icon="el-icon-lock"
            ></el-input>
          </el-form-item>
        </el-row>

        <el-row>
          <el-form-item prop="nickname" :error="signupErrNName" label="昵称：">
            <el-input
              type="text"
              v-model="signupInfo.nickname"
              placeholder="请输入一个昵称"
              auto-complete="off"
              prefix-icon="el-icon-user"
            ></el-input>
          </el-form-item>
        </el-row>

        <el-row>
          <el-form-item prop="email" :error="signupErrEml" label="邮箱：">
            <el-input
              type="text"
              v-model="signupInfo.email"
              placeholder="请输入邮箱"
              auto-complete="off"
              prefix-icon="el-icon-message"
            ></el-input>
          </el-form-item>
        </el-row>

        <el-row>
          <el-col :span="17">
            <el-form-item prop="captcha" :error="signupErrCpch" label="验证码：">
            <el-input
              type="text"
              v-model="signupInfo.captcha"
              placeholder="请输入验证码"
              auto-complete="off"
              prefix-icon="el-icon-s-claim"
            ></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="7">
            <img :src="captcha64" style="width:90%; margin-top: -5px; cursor: pointer" title="点击更换" @click="getCaptcha">
          </el-col>
        </el-row>
      </el-form>

      <el-row>
        <el-col :span="12">
          <el-button type="primary" @click="signup" style="width: 80%">注册</el-button>
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
  name: 'signup',
  created () {
    if (this.$route.query.next) {
      this.nextRoute = this.$route.query.next
    }
    this.getCaptcha()
  },
  data () {
    // eslint-disable-next-line camelcase
    let password2_validator = (rule, value, callback) => {
      if (!value) {
        callback(new Error('请确认密码'))
      } else if (value !== this.signupInfo.password) {
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
      nextRoute: {
        path: '/login'
      },
      signupInfo: {
        handle: '',
        password: '',
        password2: '',
        email: '',
        nickname: '',
        captcha: ''
      },
      captcha64: '',
      captchaKey: '',
      signupErrUsrnm: '',
      signupErrPswd: '',
      signupErrPswd2: '',
      signupErrEml: '',
      signupErrCpch: '',
      signupErrNName: '',
      rules: {
        handle: [
          {validator: this.$functions.globalValidator, trigger: 'blur'},
          {required: true, message: '请输入用户名', trigger: 'blur'},
          {min: 3, max: 30, message: '用户名应在3-30个字符之间', trigger: 'blur'}
        ],
        password: [
          {validator: this.$functions.globalValidator, trigger: 'blur'},
          {required: true, message: '请输入密码', trigger: 'blur'},
          {min: 3, max: 30, message: '密码应在3-30个字符之间', trigger: 'blur'}
        ],
        password2: [
          {validator: this.$functions.globalValidator, trigger: 'blur'},
          {required: true, message: '请再次输入密码', trigger: 'blur'},
          {validator: password2_validator, trigger: 'blur'}
        ],
        email: [
          {validator: this.$functions.globalValidator, trigger: 'blur'},
          {required: true, message: '请输入邮箱', trigger: 'blur'},
          {validator: email_validator, trigger: 'blur'}
        ],
        nickname: [
          {validator: this.$functions.globalValidator, trigger: 'blur'},
          {required: true, message: '请输入昵称', trigger: 'blur'},
          {min: 3, max: 30, message: '昵称应在3-30个字符之间', trigger: 'blur'}
        ],
        captcha: [
          {validator: this.$functions.globalValidator, trigger: 'blur'},
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
        this.captchaKey = res.data.key
      // eslint-disable-next-line handle-callback-err
      }).catch(err => {
        this.$message.error('无法获取验证码，请检查网络')
      })
    },
    signup () {
      this.$refs['signupform'].validate(valid => {
        if (valid) {
          this.signupErrUsrnm = ''
          this.signupErrPswd = ''
          this.signupErrPswd2 = ''
          this.signupErrEml = ''
          this.signupErrCpch = ''
          this.signupErrNName = ''
          const loading = this.$loading({lock: true, text: '注册中'})
          let params = this.$qs.stringify({
            'handle': this.signupInfo.handle,
            'password': this.signupInfo.password,
            'email': this.signupInfo.email,
            'nickname': this.signupInfo.nickname,
            'captcha_value': this.signupInfo.captcha,
            'captcha_key': this.captchaKey
          })
          this.$axios.post(
            '/signup',
            params
          ).then(res => {
            loading.close()
            this.$alert('注册成功，请登录', '成功', {
              confirmButtonText: '确定',
              callback: action => {
                this.$router.push({
                  path: '/login',
                  query: {
                    next: this.nextRoute
                  }
                })
              }
            })
          }).catch(err => {
            loading.close()
            this.$message.error('注册失败，请检查表单')
            err.response.data.err.forEach(item => {
              switch (parseInt(item)) {
                case 1:
                  this.signupErrUsrnm = '用户名已被注册'
                  break
                case 2:
                  this.signupErrEml = '邮箱已被注册'
                  break
                case 3:
                  this.signupErrCpch = '验证码不正确'
                  break
                default:
                  break
              }
            })
          })
        }
      })
    },
    goLogin () {
      this.$router.push({
        path: '/login',
        query: {
          next: this.nextRoute
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
