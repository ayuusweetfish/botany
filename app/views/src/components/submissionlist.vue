<template>
  <div>
    <el-collapse @change="handleCollapse">
      <el-collapse-item :name="1">
        <template slot="title">
          <div style="margin-left: 20px; font-size: 16px">手动生成对局</div>
        </template>
        <div align="left" style="margin-left: 20px">
          <el-tag
            v-for="(item, key) in selectedSubs"
            :key="key"
            closable @close="handleTagClose(item.id)"
            align="center"
            style="min-width: 100px; margin-right: 10px"
          >[{{item.id}}][<i class="el-icon-user-solid"></i>{{item.participant.handle}}]</el-tag>
          <div v-if="selectedSubs.length===0" style="color: silver">请在下方选择要进行对局的提交记录</div>
        </div>
        <div v-if="selectedSubs.length>1" align="left" style="margin-left: 20px">
          <el-link :underline="false" type="success" size="small" @click="runManualMatch">使用这些记录进行对局</el-link>
        </div>
      </el-collapse-item>
    </el-collapse>
    <div align="left" style="margin: 10px 20px 10px 20px">
      <div style="display: inline;">共有{{total}}项提交记录</div>
    </div>
    <el-table :data="submissions" v-loading="tableLoading">
      <el-table-column v-if="selecting" label="批量操作" width="80" align="center">
        <template slot-scope="scope">
          <el-button
            plain
            icon="el-icon-circle-plus-outline"
            type="text"
            size="mini"
            :disabled="!checkSelectable(scope.row) || scope.row.status!==$consts.codeStat.accepted"
            @click="addToSelections(scope.row)"
            >选择</el-button>
        </template>
      </el-table-column>
      <el-table-column label="记录编号" width="80" align="center">
        <template slot-scope="scope">
          <div>{{scope.row.id}}</div>
        </template>
      </el-table-column>
        <el-table-column label="提交人" min-width="320" align="center">
        <template slot-scope="scope">
          <div>
            <div>
              <el-avatar
              style="margin: 0;"
              :size="26"
              shape="square"
              :src="$axios.defaults.baseURL + '/user/' + scope.row.participant.handle + '/avatar'">
              </el-avatar>
            </div>
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
      <el-table-column label="语言" prop="lang" min-width="100" align="center"></el-table-column>
      <el-table-column label="提交时间" min-width="200" align="center">
        <template slot-scope="scope">
          <div>{{scope.row.timeStr}}</div>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="80" align="center">
        <template slot-scope="scope">
          <div v-if="scope.row.status===$consts.codeStat.pending" style="color: gray">等待处理</div>
          <div v-else-if="scope.row.status===$consts.codeStat.compiling" style="color: orange">处理中</div>
          <div v-else-if="scope.row.status===$consts.codeStat.accepted" style="color: green">编译通过</div>
          <div v-else-if="scope.row.status===$consts.codeStat.cmplErr" style="color: red">编译错误</div>
          <div v-else style="color: red">系统错误</div>
        </template>
      </el-table-column>
      <el-table-column label="详情" width="80" align="center">
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
      tableLoading: false,
      selectedSubs: [],
      selecting: false
    }
  },
  methods: {
    checkSelectable (row) {
      if (this.selectedSubs.find(item => (item.id === row.id || item.participant.id === row.participant.id))) {
        return false
      } else {
        return true
      }
    },
    runManualMatch () {
      if (this.selectedSubs.length < 1) {
        this.$message.warning('你需要选择至少2条记录')
        return
      }
      const loading = this.$loading({ lock: true, text: '处理中' })
      const idList = []
      this.selectedSubs.forEach(item => {
        idList.push(item.id)
      })
      const params = this.$qs.stringify({
        submissions: idList.join(',')
      })
      this.$axios.post(
        '/contest/' + this.cid + '/match/manual',
        params
      ).then(res => {
        const mid = res.data.id
        loading.close()
        this.$confirm(
          '对局生成成功',
          '提示',
          {
            type: 'success',
            confirmButtonText: '转到对局页面',
            cancelButtonText: '留在当前页面'
          }
        ).then(() => {
          this.$router.push({
            path: '/match',
            query: {
              mid: mid,
              cid: this.cid
            }
          })
        }).catch(() => {})
      }).catch(err => {
        this.$message.error('生成失败')
        loading.close()
        console.log(err)
      })
    },
    addToSelections (item) {
      this.selectedSubs.push(item)
    },
    handleCollapse (item) {
      if (item && item.length > 0) {
        this.selecting = true
      } else {
        this.selecting = false
      }
    },
    handleTagClose (val) {
      this.selectedSubs = this.selectedSubs.filter(item => item.id !== val)
    },
    handleCurrentChange (val) {
      this.page = val
      this.getList()
    },
    getList () {
      this.tableLoading = true
      const params = {
        page: this.page - 1,
        count: this.count
      }
      this.$axios.get(
        '/contest/' + this.cid + '/submission/list',
        { params: params }
      ).then(res => {
        console.log(res.data)
        this.submissions = []
        this.total = res.data.total
        res.data.submissions.forEach(item => {
          this.submissions.push({
            id: item.id,
            lang: item.lang,
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
    }
  }
}
</script>
