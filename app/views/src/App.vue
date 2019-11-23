<template>
  <div id="app">
    <el-row class = "topbar" type="flex" justify="space-between">
      <el-col :span="10" v-if="$route.meta.navbarType !== 'game'" class="topbar-tittle">
        <div align='left'>Botany-Demo</div>
      </el-col>
      <el-col :span="8" v-if="$route.meta.navbarType === 'game'" class="topbar-tittle">
        <div align='left'>GStrategy@Botany</div>
      </el-col>
      <el-col :span="6" v-if="$route.meta.navbarType === 'main'" align='left'>
        <div style="display: inline; color:gray">|</div>
        <el-button type="text" class="navbar-item" @click="goGamelist">比赛列表</el-button>
        <div style="display: inline; color:gray">|</div>
        <el-button type="text" class="navbar-item" @click="goPersonalinfo">个人信息</el-button>
        <div style="display: inline; color:gray">|</div>
      </el-col>
      <el-col :span="12" v-if="$route.meta.navbarType === 'game'" align='left'>
        <div style="display: inline; color:gray">|</div>
        <el-button type="text" class="navbar-item" @click="goGamemain">比赛首页</el-button>
        <div style="display: inline; color:gray">|</div>
        <el-button type="text" class="navbar-item" @click="goGamedetail">参赛指南</el-button>
        <div style="display: inline; color:gray">|</div>
        <el-button type="text" class="navbar-item" @click="goCoding">我的代码</el-button>
        <div style="display: inline; color:gray">|</div>
        <el-button type="text" class="navbar-item" @click="goGameranking">查看排行</el-button>
        <div style="display: inline; color:gray">|</div>
        <el-button type="text" class="navbar-item" @click="goGamevss">查看对局</el-button>
        <div style="display: inline; color:gray">|</div>
        <el-button type="text" class="navbar-item" @click="goGamelist">返回</el-button>
        <div style="display: inline; color:gray">|</div>
      </el-col>
      <el-col :span="2" v-if="$route.meta.navbarType !== 'none'">
        <el-dropdown :hide-on-click="false" @command="handleCommand" v-model="username">
          <span class="el-dropdown-link" style="cursor: pointer;">
            <el-avatar size = "small"></el-avatar>
            <p style="margin: 0">{{$store.state.username}}</p>
          </span>
          <el-dropdown-menu  slot="dropdown">
            <el-dropdown-item command="logout">退出</el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </el-col>
    </el-row>
    <el-row>
      <el-col :span="24">
        <div v-if="$route.meta.navbarType !== 'none'" align="left">
          <i class="el-icon-caret-right"></i>
          breadcrumb
        </div>
      </el-col>
    </el-row>
    <router-view/>
  </div>
</template>

<script>

export default {
  name: 'App',
  data() {
    return {
      username: '',
      showPwdDlg: false,
      showPhnDlg: false
    }
  },
  methods: {
    handleCommand(command) {
      this.$router.push('/')
    },
    goGamelist(){
      this.$router.push('/gamelist')
    },
    goGamemain(){
      this.$router.push('/gamemain')
    },
    goPersonalinfo(){
      this.$router.push('/personalinfo')
    },
    goGamedetail(){
      this.$router.push('/gamedetail')
    },
    goGameranking(){
      this.$router.push('/gameranking')
    },
    goGamevss(){
      this.$router.push('/gamevss')
    },
    goCoding(){
      this.$router.push('/coding')
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
  min-width: 720px;
  max-width: 1080px;
}
.topbar {
  border-bottom: 1px solid silver;
  max-width: 1080px;
  margin: auto;
  align-items: flex-end;
  margin-bottom: 20px;
  min-width: 720px;
  height: 84px;
}
.topbar-tittle {
  font-weight: 600;
  font-size: 30px;
  text-align: left;
  margin-left: 20px;
  line-height: 60px;
}
.navbar-item {
  font-size: 18px;
  color: grey;
  height: 30px;
}
</style>
