<template>
  <div>
    <el-row style="margin-bottom: 10px">
      <el-card body-style="height: 480px">
        <div align="left" style="font-size: 24px; font-weight: 600;">详细介绍</div>
        <div align="left" style="white-space: pre-wrap;">{{details}}</div>
      </el-card>
    </el-row>
    <el-row>
      <el-card>
        <div align="left" style="font-size: 24px; font-weight: 600">文件下载</div>
        <el-button type="primary">点击此处下载SDK</el-button>
        <el-button type="primary">点击此处下载文档</el-button>
      </el-card>
    </el-row>
  </div>
</template>

<script>
export default {
  name: 'contestdetail',
  created () {
    this.cid = this.$route.query.cid
    const loading = this.$loading({lock: true, text: '查询中'})
    this.$axios.get(
      '/contest/' + this.cid + '/info'
    ).then(res => {
      this.isRegOpen = res.data.is_reg_open
      this.myRole = res.data.my_role
      this.details = res.data.details
      this.$store.commit('enterSubSite', res.data)
      loading.close()
    }).catch(err => {
      loading.close()
      console.log(err)
      this.$message.error('查询失败')
    })
  },
  data () {
    return {
      cid: '',
      isRegOpen: false,
      myRole: -1,
      details: ''
    }
  }
}
</script>
