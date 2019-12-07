<template>
  <div>
    <el-row style="margin-bottom: 10px">
      <el-card>
        <div align="left">
          <div style="display: inline">代码已提交，状态：</div>
          <div style="display: inline; color: orange">处理中</div>
        </div>
      </el-card>
    </el-row>
    <el-row :gutter="20" style="margin-bottom: 10px">
      <el-col :span="18">
        <el-card body-style="height: 360px">
          <div align="left" style="font-size: 18px; font-weight: 600; margin-bottom: 20px">代码编辑</div>
          <codemirror
            v-model="code"
            :options="cmOptions"
            class="code"
            align="left"
          >
          </codemirror>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card body-style="height: 360px">
          <div align="left" style="font-size: 18px; font-weight: 600; margin-bottom: 20px">历史代码</div>
          <el-timeline align="left">
            <el-timeline-item
              v-for="(activity, index) in history"
              :key="index"
              :timestamp="activity.time"
              :color="activity.color"
              placement="top"
            >
            <div align="left">编号：{{activity.sid}}</div>
            <div align="left">状态：{{activity.stat}}</div>
            <el-button type="text" size="small">点击导出</el-button>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
    </el-row>
    <el-row>
      <el-card>
        <div align="left" style="font-size: 18px; font-weight: 600; margin-bottom: 20px">操作</div>
        <el-row>
          <el-col :span="8">
            <el-button type="primary" size="small" style="width: 80%">提交</el-button>
          </el-col>
          <el-col :span="8">
            <el-button size="small" style="width: 80%">保存至草稿</el-button>
          </el-col>
          <el-col :span="8">
            <el-button size="small" style="width: 80%">导出草稿</el-button>
          </el-col>
        </el-row>
      </el-card>
    </el-row>
  </div>
</template>

<script>
import {codemirror} from 'vue-codemirror-lite'

export default {
  name: 'coding',
  components: {
    codemirror
  },
  data () {
    return {
      cid: '',
      code: '',
      topbarText: '尚未提交代码',
      history: [],
      cmOptions: {
        lineNumbers: true,
        indentUnit: 2,
        autoCloseBrackets: true
      }
    }
  }
}
</script>

<style scoped>
  .code{
    border: 1px solid #dcdfe6;
    margin: auto;
  }
</style>
