<template>
  <el-container class="login-container">
    <el-header class="login-header">
      <div style="margin:auto">登录</div>
    </el-header>
    <el-main class="login-main">
      <el-form
        ref="loginform"
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
            <el-form-item prop="password" :error="loginErrPswd">
              <el-input
                type="password"
                v-model="loginInfo.password"
                placeholder="请输入密码"
                auto-complete="off"
                prefix-icon="el-icon-lock"
              ></el-input>
            </el-form-item>
          </el-col>
        </el-row>
                   
        <el-row>
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
            <img :src="captcha64" style="width:100%; margin-top: -10px" @click="getCaptcha">
          </el-col>
        </el-row>
      </el-form>

      <el-row>
        <el-col :span="12">
          <el-button type="primary" @click="login" style="width: 80%">登录</el-button>
        </el-col>
        <el-col :span="12">
          <el-button @click="goRegister" style="width: 80%">注册</el-button>
        </el-col>
      </el-row>
      <el-row>
        <el-button type="text">忘记密码？</el-button>
      </el-row>          
    </el-main>
  </el-container>
</template>

<script>
export default {
  name: 'login',
  created() {
    this.getCaptcha()
  },
  data() {
    return {
      loginInfo: {
        username: '',
        password: '',
        captcha:  '',
      },
      captcha64: '',
      loginErrUsrnm: '',
      loginErrPswd: '',
      loginErrCpch: '',
      rules: {
        username: [
          {required: true, message: '请输入账号', trigger: 'blur'},
          {max: 30, message: '输入过长', trigger:'blur'}
        ],
        password: [
          {required: true, message: '请输入密码', trigger: 'blur'},
          {max: 30, message: '输入过长', trigger:'blur'}
        ],
        captcha: [
          {required: true, message: '请输入验证码', trigger: 'blur'},
          {min: 4, max: 4, message: '请输入4个字符', trigger:'blur'}
        ],
      }
    }
  },
  methods: {
    getCaptcha() {
      this.$axios.get(
        '/captcha/login'
      ).then(res=>{
        this.captcha64 = res.data.pic
      }).catch(err=>{
        this.$message.error('无法获取验证码，请检查网络')
      })
    },
    login() {
      this.$refs['loginform'].validate(valid=>{
        if(valid){
          this.loginErrUsrnm = ''
          this.loginErrPswd = ''
          this.loginErrCpch = ''
          const loading = this.$loading({lock: true, text: '登录中'})
          let params = new URLSearchParams()
          params.append('username', this.loginInfo.username)
          params.append('password', this.loginInfo.password)
          params.append('captcha', this.loginInfo.captcha)
          this.$axios.post(
            '/login',
            params
          ).then(res=>{
            let logindata = {
              username: this.loginInfo.username,
              userid: res.data.uid,
              usertype: res.data.usertype
            }
            this.$store.commit('login', logindata)
            loading.close()
            this.$router.push('/gamelist')
          }).catch(err=>{
            loading.close()
            this.$message.error('登录失败')
            if(err.response.data.error === 'wrong captcha'){
              this.loginErrCpch = '验证码错误'
              this.loginInfo.captcha = ''
              this.getCaptcha()
            }
            else if (err.response.data.error === 'wrong username or password'){
              this.loginErrUsrnm = '用户名或密码错误'
              this.loginErrPswd = '用户名或密码错误'
              this.loginInfo.password = ''
              this.getCaptcha()
            }
          })
        }
      })
      
    },
    goRegister(){
      this.$router.push('/register')
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