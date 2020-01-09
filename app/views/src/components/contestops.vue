<template>
  <div style="margin: auto">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card shadow="hover" class="tab-cards" @click.native="goEdit">
          <div class="tab-headers" slot="header">编辑比赛</div>
          <div align="left">
            <div>比赛标题</div>
            <div>首页banner</div>
            <div>起止时间</div>
            <div>简要介绍</div>
            <div>详细介绍</div>
            <div>设置脚本</div>
            <div>设置html</div>
            <div>报名开放</div>
            <div>设置管理员</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="tab-cards" @click.native="goJudge">
          <div class="tab-headers" slot="header">裁判程序</div>
          <div align="left">
            <div>上传裁判代码</div>
            <div>选择裁判程序</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="tab-cards" @click.native="goLog">
          <div class="tab-headers" slot="header">比赛日志</div>
          <div align="left">
            <div>操作脚本</div>
            <div>查看输出</div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  name: 'contestops',
  created(){
    this.cid = this.$route.query.cid;
  },
  data(){
    return {
      cid: -1
    }
  },
  methods: {
    goEdit() {
      this.$router.push({
        path: 'contest_edit', 
        query: {cid: this.cid}
      })
    },
    goJudge() {

    },
    goLog() {
      this.$router.push({
        path: '/contest_script',
        query: {
          cid: this.cid
        }
      })
    }
  }
}
</script>

<style scoped>
  .tab-cards {
    height: 360px;
    cursor: pointer;
  }
  .tab-headers {
    font-size: 18px;
    font-weight: 600;
    text-align: left;
  }
</style>