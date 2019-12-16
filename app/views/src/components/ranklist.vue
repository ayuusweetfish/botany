<template>
  <div>
    <el-card>
      <div align="left">
        <div style="display: inline">本比赛共有{{total}}位选手</div>
      </div>
      <el-table :data="players" v-loading="tableLoading">
        <el-table-column label="排名" type="index" :index="val=>{return val + 1}" width="80"></el-table-column>
        <el-table-column label="选手" prop="participant.nickname" min-width="160" align="center">
        </el-table-column>
        <el-table-column label="账号 | ID" min-width="160" align="center">
          <template slot-scope="scope">
            <router-link
              style="display: inline; text-decoration: none; color: #409EFF"
              :to="{path: '/profile', query: {handle: scope.row.participant.handle}}"
            >{{scope.row.participant.handle}}
            </router-link>
            <el-divider direction="vertical" style="align-self: center"></el-divider>
            <div style="display: inline">{{scope.row.participant.id}}</div>
          </template>
        </el-table-column>
        <el-table-column label="评分" prop="rating" width="80" align="center">
        </el-table-column>
        <el-table-column label="表现" prop="performance" min-width="160" align="center">
        </el-table-column>
      </el-table>
      <el-pagination
        :total="total"
        :current-page="1"
        :page-size="20"
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
  name: 'ranklist',
  created () {
    this.cid = this.$route.query.cid
    this.getList()
  },
  data () {
    return {
      tableLoading: false,
      players: [],
      total: 0,
      page: 1,
      count: 20
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
        page: this.page - 1,
        count: this.count
      }
      this.$axios.get(
        '/contest/' + this.cid + '/ranklist',
        {params: params}
      ).then(res => {
        this.total = res.data.total
        this.players = res.data.participants
        this.tableLoading = false
      }).catch(err => {
        console.log(err)
        this.tableLoading = false
        this.$message.error('查询失败')
      })
    }
  }
}
</script>
