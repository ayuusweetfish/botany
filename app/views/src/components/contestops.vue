<template>
  <div style="margin: auto">
    <el-row :gutter="20" style="margin-top: 200px;">
      <el-col :span="8">
        <el-card shadow="hover" class="tab-cards" @click.native="goEdit">
          <div class="tab-headers" slot="header">编辑比赛</div>
          <div align="left">
            <div><i class="el-icon-s-tools"></i>  比赛标题</div>
            <div><i class="el-icon-s-tools"></i>  起止时间</div>
            <div><i class="el-icon-s-tools"></i>  简要介绍</div>
            <div><i class="el-icon-s-tools"></i>  详细介绍</div>
            <div><i class="el-icon-s-tools"></i>  设置脚本</div>
            <div><i class="el-icon-s-tools"></i>  播放器html</div>
            <div><i class="el-icon-s-tools"></i>  报名开放</div>
            <div><i class="el-icon-s-tools"></i>  设置管理员</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="tab-cards" @click.native="goJudge">
          <div class="tab-headers" slot="header">裁判程序</div>
          <div align="left">
            <div><i class="el-icon-edit"></i>  编辑裁判代码</div>
            <div><i class="el-icon-finished"></i>  选择裁判程序</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="tab-cards" @click.native="goLog">
          <div class="tab-headers" slot="header">赛制脚本</div>
          <div align="left">
            <div><i class="el-icon-set-up"></i>  操作脚本</div>
            <div><i class="el-icon-document"></i>  查看输出</div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  name: 'contestops',
  created () {
    this.cid = this.$route.query.cid
  },
  data () {
    return {
      cid: -1
    }
  },
  methods: {
    goEdit () {
      this.$router.push({
        path: 'contest_edit',
        query: { cid: this.cid }
      })
    },
    goJudge () {
      this.$router.push({
        path: 'judge',
        query: { cid: this.cid }
      })
    },
    goLog () {
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
    height: 280px;
    cursor: pointer;
  }
  .tab-headers {
    font-size: 18px;
    font-weight: 600;
    text-align: left;
  }
</style>
