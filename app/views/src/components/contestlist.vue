<template>
  <el-card>
    <div align="left" style="margin-left: 10px">
      <router-link
        v-if="$store.state.privilege===$consts.privilege.organizer"
        to="/contest_create"
        class="addcontest-link"
      >
      <i class="el-icon-circle-plus-outline"></i>
      添加一场比赛
      </router-link>
    </div>
    <el-table :data="contests" @row-click="goContestMain" :cell-style="{'cursor': 'pointer'}">
      <el-table-column :label="title">
        <template slot-scope="scope">
          <div>
            <div class="important">名称：</div>
            <div class="normal">{{scope.row.name}}</div>
          </div>
          <div><div class="important">时间：</div><div class="normal">{{scope.row.time}}</div></div>
          <div><div class="important">说明：</div><div class="normal">{{scope.row.info}}</div></div>
          <div v-if="scope.row.myRole===$consts.role.moderator">
            <div class="important">我管理的比赛</div>
          </div>
          <div v-else-if="scope.row.myRole===$consts.role.imIn">
            <div class="important">我参加的比赛</div>
          </div>
          <div v-if="scope.row.regOpen">
            <div class="important">开放报名中</div>
          </div>
          <div v-else><div class="normal">不开放报名</div></div>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script>
export default {
  name: 'contestlist',
  created () {
    this.getContestList()
  },
  data () {
    return {
      title: '',
      contests: []
    }
  },
  methods: {
    getContestList () {
      const loading = this.$loading({lock: true, text: '正在查询比赛列表'})
      this.$axios.get(
        '/contest/list'
      ).then(res => {
        this.contests = []
        res.data.forEach(element => {
          let timeStartStr = this.$functions.dateTimeString(element.start_time)
          let timeEndStr = this.$functions.dateTimeString(element.end_time)
          this.contests.push({
            id: element.id,
            name: element.title,
            time: timeStartStr + ' 到 ' + timeEndStr,
            info: element.desc,
            myRole: element.my_role,
            regOpen: element.is_reg_open
          })
        })
        this.total = this.contests.length
        this.title = '当前共有' + this.total + '场比赛'
        loading.close()
      // eslint-disable-next-line handle-callback-err
      }).catch(err => {
        loading.close()
        this.$message.error('查询比赛列表失败')
      })
    },
    goContestMain (obj) {
      this.$router.push({
        path: '/contest_main',
        query: {
          cid: obj.id
        }
      })
    }
  }
}
</script>

<style scoped>
  .important{
    display: inline-block;
    font-weight: 600;
  }
  .normal{
    display: inline-block;
    font-weight: 400;
  }
  .addcontest-link{
    text-decoration: none;
    color: #409EFF;
  }
</style>
