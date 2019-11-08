<template>
  <div>
    <el-card>
      <div align="left">
        <div style="display: inline">本比赛共有110位选手，你排在第</div><div style="display: inline; font-weight: 600">23</div><div style="display: inline">位</div>
      </div>
      <el-table :data="players">
        <el-table-column label="排名" min-width="80" align="center">
          <template slot-scope="scope">
            <div>{{scope.row.ranking}}</div>
          </template>
        </el-table-column>
        <el-table-column label="选手" min-width="160" align="center">
          <template slot-scope="scope">
            <el-row>
              <el-col :span="8" align="right">
                <el-avatar size="medium" style="margin-right: 10px"></el-avatar>
              </el-col>
              <el-col :span="16" align="left">
                <div>{{scope.row.name}}</div>
                <div>{{scope.row.email}}</div>
              </el-col>
            </el-row>
          </template>
        </el-table-column>
        <el-table-column label="战绩(胜场|负场|胜率)" min-width="160" align="center">
          <template slot-scope="scope">
            <div>{{scope.row.win}}|{{scope.row.loss}}|{{scope.row.rate}}</div>
          </template>
        </el-table-column>
        <el-table-column label="MMR" min-width="80" align="center">
          <template slot-scope="scope">
            <div>{{scope.row.mmr}}</div>
          </template>
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
  name: 'gameranking',
  data() {
    let ls = []
    for(let i = 0; i < 20; ++i){
      ls.push({
          ranking: (i+1).toString(),
          name: 'USERNAME',
          email: 'EMAIL@ADDR.com',
          win: '20',
          loss: '20',
          rate: '50%',
          mmr: '1000'
        })
    }
    return {
      players: ls,
      total: 110,
      page: 1
    }
  },
  methods: {
    handleCurrentChange(val){

    }
  }
}
</script>