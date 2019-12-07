<template>
  <div>
    <el-card>
      <div align="left">
        <div style="display: inline">一共进行了1000场对局，你参与了其中</div><div style="display: inline; font-weight: 600">20</div><div style="display: inline">场</div>
      </div>
      <el-table :data="contests" @row-click="goDetail">
        <el-table-column label="对局编号" min-width="80" align="center">
          <template slot-scope="scope">
            <div>{{scope.row.id}}</div>
          </template>
        </el-table-column>
        <el-table-column label="对阵" min-width="160" align="center">
          <template slot-scope="scope">
            <el-row>
              <el-col :span="4" align="right">
                <el-avatar size="small" style="margin-right: 5px"></el-avatar>
              </el-col>
              <el-col :span="7" align="left">
                <div>{{scope.row.name1}}</div>
              </el-col>
              <el-col :span="2" align="center">
                <div>VS</div>
              </el-col>
              <el-col :span="4" align="right">
                <el-avatar size="small" style="margin-right: 5px"></el-avatar>
              </el-col>
              <el-col :span="7" align="left">
                <div>{{scope.row.name2}}</div>
              </el-col>
            </el-row>
          </template>
        </el-table-column>
        <el-table-column label="胜者" min-width="100" align="center">
          <template slot-scope="scope">
            <div>{{scope.row.winner}}</div>
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
  name: 'matchlist',
  data () {
    let ls = []
    for (let i = 0; i < 20; ++i) {
      ls.push({
        id: (1000 - i).toString(),
        name1: 'USERNAME1',
        name2: 'USERNAME2',
        winner: 'USERNAME1'
      })
    }
    return {
      contests: ls,
      total: 1000,
      page: 1
    }
  },
  methods: {
    handleCurrentChange (val) {

    },
    goDetail (x, y, z) {
      this.$router.push('/vsdetail')
    }
  }
}
</script>
