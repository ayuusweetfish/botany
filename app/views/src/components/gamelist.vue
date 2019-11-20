<template>
  <el-card>
    <el-table :data="games" @row-click="goGamemain">
      <el-table-column :label="title">
        <template slot-scope="scope">
          <div><div class="important">名称：</div><div class="normal">{{scope.row.name}}</div></div>
          <div><div class="important">时间：</div><div class="normal">{{scope.row.time}}</div></div>
          <div><div class="important">说明：</div><div class="normal">{{scope.row.info}}</div></div>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script>
export default {
  name:'gamelist',
  created(){
    this.getGameList()
  },
  data() {
    return {
      title: '当前共有0场比赛正在进行',
      total: 0,
      games: []
    }
  },
  methods: {
    getGameList(){
      const loading = this.$loading({lock: true, text: '正在查询比赛列表'})
      this.$axios.get(
        '/gamelist'
      ).then(res=>{
        this.total = res.data.total
        res.data.games.forEach(element => {
          this.games.push({
            id: element.id,
            name: element.name,
            time: element.time_start + ' 到 ' + element.time_end,
            info: element.info
          })
        })
        this.title = '当前共有' + this.total + '场比赛正在进行'
        loading.close()
      }).catch(err=>{
        this.$message.error('查询比赛列表失败')
        loading.close()
      })
    },
    goGamemain(x, y, z){
      console.log('clicked')
      this.$router.push('gamemain')
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
</style>