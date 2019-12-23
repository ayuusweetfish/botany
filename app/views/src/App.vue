<template>
  <div id="app">
    <div id="topbar-background">
      <div id="navbar-container" align="left">
        <div id="navbar-banner">
          <img height="60px" :src="require('./assets/botany.png').default" @click="$router.push('/')"/>
        </div>
        <el-divider direction="vertical"></el-divider>
        <div v-if="$route.meta.navbarType === 'contest'" class="contest-title-container">{{$store.state.contestInfo.title}}</div>
        <div v-else class="contest-title-container">博坛</div>
        <div v-if="$route.meta.navbarType === 'contest'" id="router-links-container">
          <router-link class="main-router-links" :to="{path:'/contest_main', query: {cid: $route.query.cid}}">比赛首页</router-link>
          <el-dropdown
            :hide-on-click="true"
            :show-timeout="0"
            placement="bottom"
          >
            <span class="main-router-links" style="cursor: pointer;">比赛操作<i class="el-icon-arrow-down"></i></span>
            <el-dropdown-menu slot="dropdown">
              <router-link class="navbar-block" :to="{path:'/contest_detail', query: {cid: $route.query.cid}}">
                <el-dropdown-item style="width: 110px; font-size: 16px"><i class="el-icon-magic-stick"></i>参赛指南</el-dropdown-item>
              </router-link>
              <router-link v-if="checkRouteValid('imIn')&&checkTimeValid(3)" class="navbar-block" :to="{path: '/submission', query: {cid: $route.query.cid}}">
                <el-dropdown-item style="width: 110px; font-size: 16px"><i class="el-icon-cpu"></i>我的代码</el-dropdown-item>
              </router-link>
              <router-link v-if="checkRouteValid('moderator')&&checkTimeValid(3)" class="navbar-block" :to="{path: '/submission_list', query: {cid: $route.query.cid}}">
                <el-dropdown-item style="width: 110px; font-size: 16px"><i class="el-icon-document-copy"></i>提交列表</el-dropdown-item>
              </router-link>
              <router-link v-if="checkTimeValid(3)" class="navbar-block" :to="{path:'/match_list', query: {cid: $route.query.cid}}">
                <el-dropdown-item style="width: 110px; font-size: 16px"><i class="el-icon-video-play"></i>对局列表</el-dropdown-item>
              </router-link>
              <router-link v-if="checkTimeValid(3)" class="navbar-block" :to="{path:'/ranklist', query: {cid: $route.query.cid}}">
                <el-dropdown-item style="width: 110px; font-size: 16px"><i class="el-icon-trophy"></i>选手排行</el-dropdown-item>
              </router-link>
            </el-dropdown-menu>
          </el-dropdown>
          <router-link class="main-router-links" to="/">返回Botany</router-link>
        </div>
        <div v-if="$route.meta.navbarType!=='none'" id="avatar-container">
          <el-dropdown v-if="$store.state.handle" :hide-on-click="true" @command="handleCommand" trigger="click">
            <span class="el-dropdown-link" style="cursor: pointer;">
              <el-avatar :src="defaultAva" size="medium" style="margin-top: 4px"></el-avatar>
            </span>
            <el-dropdown-menu  slot="dropdown" style="min-width: 120px; padding: 2px 10px 2px 10px;">
              <el-dropdown-item :disabled="true" class="info-dropdown-item" style="color: #505050; font-weight: 600; border-bottom: 1px solid silver">{{$store.state.nickname}}</el-dropdown-item>
              <!-- <el-dropdown-item :disabled="true" class="info-dropdown-item">UID: {{$store.state.id}}</el-dropdown-item> -->
              <el-dropdown-item :disabled="true" class="info-dropdown-item" style="color: grey">{{translatePrivilege($store.state.privilege)}}</el-dropdown-item>
              <el-dropdown-item :disabled="true" class="info-dropdown-item">账号：{{$store.state.handle}}</el-dropdown-item>
              <el-dropdown-item :disabled="true" class="info-dropdown-item">UID：{{$store.state.id}}</el-dropdown-item>
              <router-link :to="{path: '/profile', query: {handle: $store.state.handle}}" style="text-decoration: none">
                <el-dropdown-item class="button-dropdown-item" style="border-top: 1px solid silver">我的资料</el-dropdown-item>
              </router-link>
              <el-dropdown-item command="password" class="button-dropdown-item">修改密码</el-dropdown-item>
              <el-dropdown-item command="logout" class="button-dropdown-item">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
          <div v-else id="login-button-container">
            <el-link :underline="false" type="primary" class="login-button" style="font-size: 18px" @click="goLogin">登录</el-link>
            <el-divider direction="vertical"></el-divider>
            <el-link :underline="false" class="login-button" style="font-size: 18px" @click="goSignup">注册</el-link>
          </div>
        </div>
      </div>
      <div v-if="$route.meta.navbarType !== 'none'" align="left" id="breadcrumb-container">
        <div v-if="$store.state.routeList.length !== 0" id="breadcrumb-start">location =</div>
        <div v-for="(item, index) in $store.state.routeList" :key="index" style="display: inline">
          <router-link v-if="index!==$store.state.routeList.length-1" :to="{path: item.path, query: item.query}" class="breadcrumb-item">{{item.title}}</router-link>
          <div v-if="index!==$store.state.routeList.length-1" class="breadcrumb-connector">-></div>
          <div v-else id="breadcrumb-tail">{{item.title}}</div>
        </div>
      </div>
    </div>
    <password-dialog :visible.sync="showPwdDlg" @setVisible="setPasswordDialog"></password-dialog>
    <router-view class="main-view" style="margin-top: 100px"/>
  </div>
</template>

<script>
import passwordDialog from './components/password.vue'
export default {
  name: 'App',
  components: {
    'password-dialog': passwordDialog
  },
  created () {
    window.onscroll = function () {
      const sl = -Math.max(document.body.scrollLeft, document.documentElement.scrollLeft)
      document.getElementById('topbar-background').style.left = sl + 'px'
      const st = Math.max(document.body.scrollTop, document.documentElement.scrollTop)
      if (st > 70) {
        const top = document.getElementById('topbar-background')
        top.style.borderBottom = '1px solid #dddddd'
      } else {
        const top = document.getElementById('topbar-background')
        top.style.borderBottom = 'none'
      }
    }
    this.$axios.get('/whoami').then(res => {
      console.log(res.data)
      this.$store.commit('login', {
        id: res.data.id,
        handle: res.data.handle,
        privilege: res.data.privilege,
        nickname: res.data.nickname
      })
    }).catch(err => {
      if (err.response.state === 400) {
        this.$store.commit('logout')
      }
    })
  },
  beforeDestroy () {
    this.$store.commit('logout')
  },
  data () {
    return {
      showPwdDlg: false,
      defaultAva: require('@/assets/logo.png').default
    }
  },
  methods: {
    changeTitle (title) {
      console.log(title)
    },
    setPasswordDialog (val) {
      this.showPwdDlg = val
    },
    goLogin () {
      this.$store.commit('setAfterLogin', { path: this.$route.path, query: this.$route.query })
      this.$router.push({
        path: '/login',
        query: {
          redirect: (this.$route.path !== '/notfound')
        }
      })
    },
    goSignup () {
      this.$store.commit('setAfterLogin', { path: this.$route.path, query: this.$route.query })
      this.$router.push({
        path: '/signup',
        query: {
          redirect: (this.$route.path !== '/notfound')
        }
      })
    },
    handleCommand (cmd) {
      if (cmd === 'logout') {
        this.logout()
      } else if (cmd === 'toProfile') {
        this.toProfile()
      } else if (cmd === 'password') {
        this.showPwdDlg = true
      }
    },
    translatePrivilege (num) {
      let str = ''
      switch (num) {
        case this.$consts.privilege.common:
          str = 'Common User'
          break
        case this.$consts.privilege.oranizer:
          str = 'Organizer'
          break
        case this.$consts.privilege.superuser:
          str = 'Super User'
          break
        default:
          break
      }
      return str
    },
    logout () {
      this.$axios.post('/logout').then(res => {
        this.$store.commit('logout')
        window.location.reload()
      }).catch(err => {
        console.log(err)
      })
    },
    checkRouteValid (role) {
      if (this.$store.state.contestInfo &&
      this.$store.state.contestInfo.my_role !== null &&
      this.$store.state.contestInfo.my_role !== undefined &&
      this.$store.state.contestInfo.my_role === this.$consts.role[role]) {
        return true
      } else {
        return false
      }
    },
    checkTimeValid (val) {
      if (!this.$store.state.contestInfo) {
        return false
      }
      const start = this.$store.state.contestInfo.start_time
      const end = this.$store.state.contestInfo.end_time
      const stage = this.$functions.checkTime(start, end)
      if (stage === this.$consts.contestStat.going) {
        return true
      } else if (stage === val) {
        return true
      }
      return false
    }
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin: auto;
  margin-top: 0px;
  font-family:'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}
#topbar-background {
  position: fixed;
  background-color: white;
  z-index: 90;
  height: 86px;
  top: 0px;
  left: 0px;
  width: 100%;
  min-width: 1080px;
  overflow: hidden;
}
#navbar-container {
  background-color: white;
  z-index: 100;
  width: 1080px;
  margin: auto;
  top: 0px;
  line-height: 44px;
  height: 60px;
}
#navbar-banner {
  vertical-align: middle;
  display: inline-block;
  cursor: pointer;
}
.contest-title-container{
  font-size: 24px;
  font-weight: 600;
  display: inline-block;
  align-self: auto;
  padding: 0
}
#contest-title{
  font-size: 24px;
  font-weight: 600;
  display: inline-block;
}
#router-links-container{
  display: inline-block;
  margin-left: 20px;
}
#avatar-container{
  margin-top: 10px;
  margin-right: 8px;
  float: right;
}
#login-button-container{
  float: right;
}
#my-nickname{
  display: inline;
  font-size: 14px;
}
#breadcrumb-container{
  width: 1080px;
  margin: auto;
  align-items: center

}
#breadcrumb-start{
  margin-left: 5px;
  margin-right: 5px;
  display: inline-block;
  font-size: 14px;
}
.breadcrumb-item{
  display: inline;
  text-decoration: none;
  color: #555555;
  font-size: 14px;
  font-weight: 600;
}
.breadcrumb-connector{
  color: gray;
  display: inline;
  font-size: 14px;
}
#breadcrumb-tail{
  color: black;
  display: inline;
  font-size: 14px;
}
.main-view{
  width: 1080px;
  margin: auto;
}
.main-router-links{
  font-size: 18px;
  text-decoration: none;
  color: gray;
  margin-right: 20px
}
.topbar {
  border-bottom: 1px solid silver;
  max-width: 1080px;
  margin: auto;
  align-items: flex-end;
  margin-bottom: 10px;
  min-width: 720px;
}
.topbar-tittle {
  font-weight: 600;
  font-size: 30px;
  text-align: left;
  margin-left: 20px;
}
.navbar-item {
  font-size: 18px;
  font-weight: 500;
  color:#545454;
  text-decoration: none;
  margin: 0px 10px 0px 10px;
}
.navbar-block {
  font-weight: 500;
  color:#545454;
  text-decoration: none;
  display: block;
}
.info-dropdown-item {
  font-size: 14px;
  padding: 0
}
.button-dropdown-item {
  font-size: 14px;
  padding: 0;
}
.login-button {
  font-size: 18px;
  font-weight: 600;
}
.cm-container .CodeMirror {
  height: auto;
  max-height: 480px;
  font-family: monospace;
  position: relative;
  background: white;
  direction: ltr;
  overflow: hidden;
}
.cm-container .CodeMirror-wrap {
  height: auto;
  font-family: monospace;
  position: relative;
  background: white;
  direction: ltr;
}
.cm-container .CodeMirror-scroll {
  min-height: 300px;
  max-height: 480px;
}
</style>
