
<template>
  <div>
    <div class="wrap">
      <div class="svg-wrap">
        <svg></svg>
      </div>
    </div>
  </div>
</template>
<script>
import * as d3 from 'd3'

export default {
  data() {
    return {
      height: 150,
      width: 0,
      svg: {},
      innerWrap: {},
      outterWrap: {},
      innerSvg: {},
      outterSvg: {}
    }
  },
  mounted() {
    this.steps.forEach((item, index) => {
      item.index = index
    })
    this.width = Number(d3.select('.svg-wrap').style('width').replace(/px/g, ''))
    this.svg = d3.select('.svg-wrap svg')
      .attr('width', this.width)
      .attr('height', this.height)
      .append('g')
      .attr('transform', `translate(${-(this.width + 300)}, 0)`)
    this.innerWrap = this.svg.append('g')
    this.outterWrap = this.svg.append('g')
    this.render()
    this.addResize()
  },
  methods: {
    addResize() {
      window.onresize = () => {
        this.$nextTick(() => {
          this.width = Number(d3.select('.svg-wrap').style('width').replace(/px/g, ''))
          this.svg.attr('width', this.width).attr('height', this.height)
          this.updateSvg()
        })
      }
    },
    render() {
      this.outterSvg = this.setData({svg: this.outterWrap, data: this.outterSteps})
      this.innerSvg = this.setData({svg: this.innerWrap, data: this.innerSteps})
      this.updateSvg()
    },
    setData(params = {}) {
      return params.svg
        .selectAll('g')
        .data(params.data)
        .enter()
        .append('g')
        .attr('transform', '')
        .on('click', (e, d) => {
          if (d.selected) {
            this.$emit('change', d.id)
          }
          if (d.index <= this.currentStep && !d.parentId) {
            this.$emit('change', d.index)
          }
        })
    },
    updateSvg(params = {}) {
      this.renderStep({
        svg: this.innerSvg,
        data: this.innerSteps
      })
      this.renderStep({
        svg: this.outterSvg,
        data: this.outterSteps
      })
      this.svg.transition()
        .duration(1000)
        .delay(30)
        .attr('transform', `translate(0, 0)`)
    },
    renderStep(params = {}) {
      params.svg.html('')
      params.svg.append('text')
        .text(function (d) {
          return d.title
        })
        .attr('y', 55)
        .attr('x', 15)
        .attr('class', (d, i) => i === this.current ? 'selected' : '')
      params.svg.append('path')
        .attr('class', (d, i) => {
          let className = ''
          if (typeof d.parentId !== 'undefined') {
            className = this.current === d.id ? 'path finish' : 'path'
          } else {
            className = i < this.current ? 'path finish' : 'path'
          }
          return className
        })
        .attr('d', (d, i) => {
          let pathStr = `M35 17 L${this.itemWidth - 3} 17 Z`
          if (typeof d.parentIndex !== 'undefined') {
            // 分叉的路径显示
            const his = -75 * i + 57
            pathStr = `M${-this.itemWidth + 32} ${his} L35 17 Z`
          } else {
            if (i === params.data.length - 1) pathStr = ''
          }
          return pathStr
        })
      params.svg.attr('class', (d, i) => {
        let className = 'progress-item '
        if (typeof d.parentIndex !== 'undefined') {
          className += d.selected ? 'finish ' : ''
          className += this.current === d.id ? 'selected' : ''
        } else {
          className += i <= this.currentStep ? 'finish ' : ''
          className += d.selected ? 'finish ' : ''
          className += i === this.current ? 'selected' : ''
        }
        return className
      })
      params.svg.append('circle')
        .attr('r', 19)
        .attr('cx', 16)
        .attr('cy', 15)
      params.svg.append('g').html((d) => `${d.svg}`).attr('class', 'icon')
      params.svg.transition()
        .attr('transform', (d, i) => {
          let num = i
          let his = 50
          if (typeof d.parentId !== 'undefined') {
            num = d.parentIndex
            his = 75 * i + 10
          }
          return `translate(${num * (this.itemWidth) + 15}, ${his})`
        })
    }
  },
  destroyed() {
    window.onresize = null
  },
  watch: {
    current() {
      this.updateSvg()
    },
    currentStep() {
      this.updateSvg()
    }
  },
  computed: {
    plus() {
      return (this.itemWidth - 35) / this.steps.length
    },
    itemWidth() {
      return (this.width - 58) / (this.steps.length - 1)
    },
    outterSteps() {
      const arr = this.steps.filter(item => !item.child)
      return arr
    },
    innerSteps() {
      let rst = []
      this.steps.forEach((item, index) => {
        if (item.child) {
          rst = item.child.map(step => {
            step.parentIndex = index
            return step
          })
        }
      })
      return rst
    }
  },
  props: {
    current: {
      type: Number,
      default: 0
    },
    currentStep: {
      type: Number,
      default: 0
    },
    steps: {
      type: Array,
      default: () => {
        return []
      }
    }
  }
}
</script>

<style lang="less" scoped>

</style>
