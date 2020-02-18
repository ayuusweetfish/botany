<template>
  <div  align="left">
    <div style="display: inline; margin-left: 10px;">{{title}}</div>
    <div style="margin-left: 10px; display: inline">
      <router-link
        v-if="$store.state.privilege===$consts.privilege.organizer"
        to="/contest_create"
        class="addcontest-link"
      >
      <i class="el-icon-circle-plus-outline"></i>
      添加一场比赛
      </router-link>
    </div>
    <el-row :gutter="20">
      <el-col
        :span="12"
        v-for="(item, index) in contests"
        :key="index"
      >
        <el-card
          shadow="never"
          class="contest-card"
          >
          <div slot="header">
            <span class="important" @click="goContestMain(item)" style="cursor: pointer">{{item.name}}</span>
            <div v-if="item.myRole===$consts.role.moderator" style="float: right">我管理的比赛</div>
            <div v-else-if="item.myRole===$consts.role.imIn" style="float: right">我参加的比赛</div>
          </div>
          <div class="contest-info">
            <div>
              <div class="important">时间：</div>
              <div class="normal">{{item.time}}</div>
            </div>
            <div>
              <div class="important">简介：</div>
              <div class="normal">{{item.info}}</div>
            </div>
          </div>
          <div align="right">
            <router-link
              :to="{path: '/contest_main', query:{cid: item.id}}"
              class="contest-button">查看详情</router-link>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
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
      contests: [],
      listLoading: false
    }
  },
  watch: {
    '$store.state.id' () {
      this.getContestList()
    }
  },
  methods: {
    getContestList () {
      this.listLoading = true
      this.$axios.get(
        '/contest/list'
      ).then(res => {
        this.contests = []
        res.data.forEach(element => {
          const timeStartStr = this.$functions.dateTimeString(element.start_time)
          const timeEndStr = this.$functions.dateTimeString(element.end_time)
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
        this.listLoading = false
      // eslint-disable-next-line handle-callback-err
      }).catch(err => {
        this.listLoading = false
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
  .contest-card{
    text-align: left;
    margin-top: 30px;
  }
  .contest-info{
    font-size: 14px;
    height: 120px;
  }
  .contest-button{
    font-size: 14px;
    text-decoration: none;
    color: #409EFF;
  }
</style>
