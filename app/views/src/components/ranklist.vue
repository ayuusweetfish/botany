<template>
  <div>
    <div align="left">
      <div style="display: inline">本比赛共有{{total}}位选手</div>
    </div>
    <el-table :data="players" v-loading="tableLoading">
      <el-table-column label="排名" type="index" :index="val=>{return val + 1}" width="80" align="center"></el-table-column>
      <el-table-column  width="80" align="right">
        <template slot-scope="scope">
          <div style="margin-top: 6px">
            <el-avatar
            style="margin: 0;"
            :size="32"
            shape="square"
            :src="$axios.defaults.baseURL + '/user/' + scope.row.participant.handle + '/avatar'">
            </el-avatar>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="选手" prop="participant.nickname" min-width="100" align="left">
      </el-table-column>
      <el-table-column label="账号" min-width="100" align="center">
        <template slot-scope="scope">
          <router-link
            style="display: inline; text-decoration: none; color: #409EFF"
            :to="{path: '/profile', query: {handle: scope.row.participant.handle}}"
          >{{scope.row.participant.handle}}
          </router-link>
        </template>
      </el-table-column>
      <el-table-column label="主战提交" align="center">
        <template slot-scope="scope">
          <div v-if="scope.row.delegate===-1" style="color: silver">暂无</div>
          <div v-else>
            <router-link
              v-if="myRole===$consts.role.moderator"
              style="display: inline; text-decoration: none; color: #409EFF"
              :to="{path: '/submission_info', query: {cid: cid, sid: scope.row.delegate}}"
              >{{scope.row.delegate}}</router-link>
            <div v-if="myRole!==$consts.role.moderator">{{scope.row.delegate}}</div>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="评分" prop="rating" width="80" align="center">
      </el-table-column>
      <el-table-column label="表现" prop="performance" min-width="200" align="center">
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
  </div>
</template>

<script>
export default {
  name: 'ranklist',
  created () {
    this.cid = this.$route.query.cid
    this.getInfo()
    this.getList()
  },
  data () {
    return {
      myRole: -1,
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
    getInfo () {
      this.$axios.get(
        '/contest/' + this.cid + '/info'
      ).then(res => {
        this.myRole = res.data.my_role
        this.$store.commit('enterSubSite', res.data)
      }).catch(err => {
        console.log(err)
      })
    },
    getList () {
      this.tableLoading = true
      const params = {
        page: this.page - 1,
        count: this.count
      }
      this.$axios.get(
        '/contest/' + this.cid + '/ranklist',
        { params: params }
      ).then(res => {
        this.total = res.data.total
        this.players = res.data.participants
        console.log(this.players)
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
