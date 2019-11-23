<template>
  <div>
    <el-row style="margin-top: 20px; margin-bottom: 20px">
      <el-col :span="10">
        <el-avatar :size="80"></el-avatar>
        <div style="font-size: 22; font-weight: 600">USERNAME1</div>
        <div style="font-size: 22;">USER1@mails.com</div>
        <div style="font-size: 22; font-weight: 600; color: green">WINNER</div>
      </el-col>
      <el-col :span="4">
        <div style="font-size: 36px; font-weight: 600; margin-top: 36px">VS</div>
      </el-col>
      <el-col :span="10">
        <el-avatar :size="80"></el-avatar>
        <div style="font-size: 22; font-weight: 600">USERNAME2</div>
        <div style="font-size: 22;">USER2@mails.com</div>
        <div style="font-size: 22; font-weight: 600">LOSER</div>
      </el-col>
    </el-row>
    <el-row>
      <el-card body-style="height: 480px">
        <div align="left" style="font-size: 24px; font-weight: 600">比赛回放</div>
      </el-card>
    </el-row>
    <el-row>
      <el-card>
        <el-row>
          <el-col :span="6">
            <el-button style="width: 80%">上一回合</el-button>
          </el-col>
          <el-col :span="6">
            <el-button style="width: 80%" type="primary">暂停/开始</el-button>
          </el-col>
          <el-col :span="6">
            <el-button style="width: 80%">下一回合</el-button>
          </el-col>
          <el-col :span="6">
            <el-button style="width: 80%" type="text">下载录像文件</el-button>
          </el-col>
        </el-row>
      </el-card>
    </el-row>
  </div>
</template>

<script>
export default {
  name: 'vsdetail'
}
</script>
