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
  data() {
    return {
      title: '当前共有2场比赛正在进行',
      games: [
        {
          name: 'GStrategy',
          time: '2019-09-01 to 2020-01-01',
          info: 'a SSAST game.'
        },
        {
          name: 'GStrategy2',
          time: '2019-09-02 to 2020-01-02',
          info: 'a SSAST game.'
        },
      ]
    }
  },
  methods: {
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