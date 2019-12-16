<template>
  <div>
    <el-card>
      <div align="left">
        <div style="display: inline">共有{{total}}项提交记录</div>
      </div>
      <el-table :data="submissions" v-loading="tableLoading">
        <el-table-column label="记录编号" width="80" align="center">
          <template slot-scope="scope">
            <div>{{scope.row.id}}</div>
          </template>
        </el-table-column>
          <el-table-column label="提交人" width="360" align="center">
          <template slot-scope="scope">
            <div>
              <div style="display: inline; font-weight: 600">{{scope.row.participant.nickname}}</div>
              <el-divider direction="vertical"></el-divider>
              <router-link
                style="display: inline; color: #409EFF; text-decoration: none"
                :to="{path: '/profile', query: {handle: scope.row.participant.handle}}"
              >
                {{scope.row.participant.handle}}
              </router-link>
            </div>
          </template>
        </el-table-column>
          <el-table-column label="提交时间" width="280" align="center">
          <template slot-scope="scope">
            <div>{{scope.row.timeStr}}</div>
          </template>
        </el-table-column>
        <el-table-column label="状态" min-width="100" align="center">
          <template slot-scope="scope">
            <div v-if="scope.row.status===$consts.codeStat.pending" style="color: gray">等待处理</div>
            <div v-else-if="scope.row.status===$consts.codeStat.compiling" style="color: orange">处理中</div>
            <div v-else-if="scope.row.status===$consts.codeStat.accepted" style="color: accepted">已结束</div>
            <div v-else-if="scope.row.status===$consts.codeStat.cmplErr" style="color: red">编译错误</div>
            <div v-else style="color: red">Fatal Error</div>
          </template>
        </el-table-column>
        <el-table-column label="详情" min-width="100" align="center">
          <template slot-scope="scope">
            <router-link
              style="color: #409EFF; text-decoration: none"
              :to="{path: '/submission_info', query: {cid: cid, sid: scope.row.id}}"
            >查看提交</router-link>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        :total="total"
        :current-page="page"
        :page-size="count"
        @current-change="handleCurrentChange"
        :pager-count="11"
        layout="prev, pager, next, jumper, ->, total"
      >
      </el-pagination>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'submissionlist',
  created () {
    this.cid = this.$route.query.cid
    this.getList()
  },
  data () {
    return {
      cid: 0,
      submissions: [],
      total: 0,
      page: 1,
      count: 10,
      tableLoading: false
    }
  },
  methods: {
    handleCurrentChange (val) {
      this.page = val
      this.getList()
    },
    getList () {
      this.tableLoading = true
      let params = {
        'page': this.page - 1,
        'count': this.count
      }
      this.$axios.get(
        '/contest/' + this.cid + '/submission/list',
        {params: params}
      ).then(res => {
        console.log(res.data)
        this.submissions = []
        this.total = res.data.total
        res.data.submissions.forEach(item => {
          this.submissions.push({
            id: item.id,
            participant: item.participant,
            timeStr: this.$functions.dateTimeString(item.created_at),
            status: item.status
          })
        })
        this.tableLoading = false
      }).catch(err => {
        console.log(err)
        this.$message.error('查询失败')
        this.tableLoading = false
        this.submissions = []
      })
    },
    goDetail (x, y, z) {
      
    }
  }
}
</script>
