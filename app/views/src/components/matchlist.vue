<template>
  <div>
    <el-card>
      <div align="left">
        <div style="display: inline">共进行了{{total}}场对局</div>
      </div>
      <el-table :data="matches" v-loading="tableLoading">
        <el-table-column label="对局编号" width="120" align="center">
          <template slot-scope="scope">
            <div>{{scope.row.id}}</div>
          </template>
        </el-table-column>
          <el-table-column label="参赛者" align="center">
          <template slot-scope="scope">
            <div>
              <el-popover v-for="(item, index) in scope.row.parties" :key="index" style="margin: 5px" placement="bottom">
                <el-tag slot="reference" type="primary" style="cursor: pointer; min-width: 120px; text-align: center">{{item.participant.nickname}}</el-tag>
                <div>
                  <div>
                    <div class="party-item-title">代码ID:</div>
                    <div class="party-item-text">{{item.id}}</div>
                  </div>
                  <div class="party-item-title" style="display: block">选手:</div>
                  <div>
                    <div class="party-item-light">昵称</div>
                    <div class="party-item-text">{{item.participant.nickname}}</div>
                  </div>
                  <div>
                    <div class="party-item-light">主页</div>
                    <router-link class="party-item-button" :to="{path: '/profile', query: {handle: item.participant.handle}}">{{item.participant.handle}}</router-link>
                  </div>
                </div>
              </el-popover>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="160" align="center">
          <!-- <template slot-scope="scope">
            <div style="color: green">{{scope.row.res}}</div>
          </template> -->
          <template slot-scope="scope">
            <div v-if="scope.row.status===$consts.codeStat.pending" style="color: gray">等待处理</div>
            <div v-else-if="scope.row.status===$consts.codeStat.compiling" style="color: orange">处理中</div>
            <div v-else-if="scope.row.status===$consts.codeStat.accepted" style="color: green">已结束</div>
            <div v-else style="color: red">系统错误</div>
          </template>
        </el-table-column>
        <el-table-column label="选项" width="160" align="center">
          <!-- <template slot-scope="scope">
            <div style="color: green">{{scope.row.res}}</div>
          </template> -->
          <template slot-scope="scope">
            <router-link
              style="color: #409EFF; text-decoration: none"
              :to="{path: '/match', query: {cid: scope.row.contest.id, mid: scope.row.id}}"
            >查看详情
            </router-link>
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
  name: 'matchlist',
  created () {
    this.cid = this.$route.query.cid
    this.getList()
  },
  data () {
    return {
      cid: 0,
      matches: [],
      total: 0,
      page: 1,
      count: 20,
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
        '/contest/' + this.cid + '/matches',
        {params: params}
      ).then(res => {
        console.log(res.data)
        this.matches = res.data.matches
        this.total = res.data.total
        this.tableLoading = false
      }).catch(err => {
        console.log(err)
        this.$message.error('查询失败')
        this.tableLoading = false
        this.matches = []
      })
    },
    goDetail (x, y, z) {

    }
  }
}
</script>

<style scoped>
.party-item-title{
  color: gray;
  font-weight: 600;
  display: inline;
}
.party-item-light{
  color: gray;
  font-weight: 400;
  display: inline;
  margin-left: 8px;
}
.party-item-text{
  color: black;
  font-weight: 400;
  display: inline;
}
.party-item-button{
  text-decoration: none;
  color: #409EFF;
  font-weight: 400;
  display: inline;
}
</style>
