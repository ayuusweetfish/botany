<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="7">
        <el-card shadow="hover" body-style="margin-bottom: 20px">
          <div style="font-size: 18px; font-weight: 600; margin-bottom: 20px; margin-top: 20px;">我@Botany</div>
          <el-avatar :size="size" style="margin-bottom: 20px"></el-avatar>
          <div style="font-size: 22; word-wrap: break-word;">USER</div>
          <div style="font-size: 22; word-wrap: break-word;">USER@mails.tsinghua.edu.cn</div>
          <div style="font-size: 22; word-wrap: break-word;">USER.github.com</div>
          <el-button type="text">修改个人信息</el-button>
          <div style="font-size: 18px; font-weight: 600; margin-bottom: 10px; margin-top: 40px">成就</div>
          <div style="font-size: 22; word-wrap: break-word;">第一届GStrategy大赛冠军</div>
        </el-card>
      </el-col>
      <el-col :span="17">
        <el-card  style="margin-bottom: 20px">
          <div align="left">
            <div style="display: inline">一共参加了</div><div style="display: inline; font-weight: 600">2</div><div style="display: inline">项赛事</div>
          </div>
          <el-table :data="major">
            <el-table-column label="比赛" min-width="160" align="center">
              <template slot-scope="scope">
                <div>{{scope.row.name}}</div>
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
            <el-table-column label="排名" min-width="80" align="center">
              <template slot-scope="scope">
                <div>{{scope.row.rank}}</div>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
        <el-card>
          <div align="left">
            <div style="display: inline">一共进行了</div><div style="display: inline; font-weight: 600">123</div><div style="display: inline">场对局</div>
          </div>
          <el-table :data="minor">
            <el-table-column label="对局编号" min-width="80" align="center">
              <template slot-scope="scope">
                <div>{{scope.row.id}}</div>
              </template>
            </el-table-column>
             <el-table-column label="对手" min-width="160" align="center">
              <template slot-scope="scope">
                <el-row>
                  <el-col :span="12" align="right">
                    <el-avatar size="small" style="margin-right: 10px"></el-avatar>
                  </el-col>
                  <el-col :span="12" align="left">
                    <div>{{scope.row.oppo}}</div>
                  </el-col>
                </el-row>
              </template>
            </el-table-column>
            <el-table-column label="结果" min-width="80" align="center">
              <template slot-scope="scope">
                <div style="color: green">{{scope.row.res}}</div>
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
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  name: 'profile',
  data () {
    return {
      major: [
        {
          name: 'GStrategy',
          win: '20',
          loss: '20',
          rate: '50%',
          mmr: '1000',
          rank: '23'
        },
        {
          name: 'GStrategy2',
          win: '20',
          loss: '20',
          rate: '50%',
          mmr: '1000',
          rank: '23'
        }
      ],
      minor: [
        {
          id: '987',
          game: 'GStrategy',
          oppo: 'USER1',
          res: 'WON'
        },
        {
          id: '986',
          game: 'GStrategy2',
          oppo: 'USER2',
          res: 'WON'
        },
        {
          id: '985',
          game: 'GStrateg2',
          oppo: 'USER1',
          res: 'WON'
        },
        {
          id: '984',
          game: 'GStrateg2',
          oppo: 'USER1',
          res: 'WON'
        },
        {
          id: '982',
          game: 'GStrategy',
          oppo: 'USER1',
          res: 'WON'
        }
      ],
      total: 123,
      size: 120
    }
  },
  methods: {
    handleCurrentChange (val) {

    }
  }
}
</script>
