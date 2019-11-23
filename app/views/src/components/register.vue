<template>
  <el-container class="register-container">
    <el-header class="register-header">
      <div style="margin:auto">注册</div>
    </el-header>
    <el-main class="register-main">
      <el-form
        ref="regisform"
        :model="regisInfo"
        label-suffix="left"
        label-width="0px"
        :rules="rules"
      >
        <el-row>
          <el-col :span="5">
            <div align="right" class="register-title">用户名：</div>
          </el-col>
          <el-col :span="19">
            <el-form-item prop="username" :error="regisErrUsrnm">
              <el-input
                type="text"
                v-model="regisInfo.username"
                placeholder="请输入账号"
                auto-complete="off"
                prefix-icon="el-icon-user-solid"
              ></el-input>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row>
          <el-col :span="5">
            <div align="right" class="register-title">密码：</div>
          </el-col>
          <el-col :span="19">
            <el-form-item prop="password" :error="regisErrPswd">
              <el-input
                type="password"
                v-model="regisInfo.password"
                placeholder="请输入密码"
                auto-complete="off"
                prefix-icon="el-icon-lock"
              ></el-input>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row>
          <el-col :span="5">
            <div align="right" class="register-title">确认密码：</div>
          </el-col>
          <el-col :span="19">
            <el-form-item prop="password2" :error="regisErrPswd2">
              <el-input
                type="password"
                v-model="regisInfo.password2"
                placeholder="请再次输入密码"
                auto-complete="off"
                prefix-icon="el-icon-lock"
              ></el-input>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row>
          <el-col :span="5">
            <div align="right" class="register-title">邮箱：</div>
          </el-col>
          <el-col :span="19">
            <el-form-item prop="email" :error="regisErrEml">
              <el-input
                type="text"
                v-model="regisInfo.email"
                placeholder="请输入邮箱"
                auto-complete="off"
                prefix-icon="el-icon-message"
              ></el-input>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row>
          <el-col :span="5">
            <div align="right" class="register-title">验证码：</div>
          </el-col>
          <el-col :span="12">
            <el-form-item prop="captcha" :error="regisErrCpch">
            <el-input
              type="text"
              v-model="regisInfo.captcha"
              placeholder="请输入验证码"
              auto-complete="off"
              prefix-icon="el-icon-s-claim"
            ></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="7">
            <img :src="captcha64" style="width:100%; margin-top: -10px">
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
  name: 'register',
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
      // eslint-disable-next-line no-useless-escape
      } else if (!/^([a-zA-Z0-9]+[-_\.]?)+@([a-zA-Z0-9]+\.)+[a-z]+$/.test(value)) {
        callback(new Error('请输入格式正确的邮箱'))
      } else {
        callback()
      }
    }
    return {
      regisInfo: {
        username: '',
        password: '',
        password2: '',
        email: '',
        captcha: ''
      },
      captcha64: '',
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
        '/captcha/register'
      ).then(res => {
        this.captcha64 = res.data.pic
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
.register-container {
  border-radius: 5px;
  background-clip: padding-box;
  margin: 100px auto;
  padding: 20px;
  width: 550px;
  border: 1px solid silver;
  font-weight: 600;
  font-size: 25px;
}
.register-header {
  display: flex;
  text-align: center;
  border-top: none;
  border-left: none;
  border-right: none;
  border-bottom: 1px solid silver;
}
.register-main {
  font-weight: 400;
  font-size: 18px;
}
.register-title {
  font-weight: 400;
  font-size: 18px;
  vertical-align: middle;
  line-height: 40px;
}
.register-tip {
  font-weight: 400;
  font-size: 16px;
  color: silver;
  margin-bottom: 0;
  padding: 0;
  line-height: 22px;
}
</style>
