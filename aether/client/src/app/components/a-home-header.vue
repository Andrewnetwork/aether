<template>
  <div class="view-header">
    <!-- <div class="side-container left"></div> -->
    <div class="center-container">
      <div class="header-container side-layout">
        <div class="header-text">Home</div>
        <div class="header-subtext">Most popular content from your subbed communities
          <div class="infomark-container">
            <a-info-marker text="<p>If you've subbed to a community, and you have its notifications enabled, most interesting stuff in it will appear here.</p>"></a-info-marker>
          </div>
        </div>
      </div>
    </div>
    <div class="side-container right">
      <div class="actions-container">
        <router-link to="/globalscope/subbed" class="button is-link is-outlined" title="Edit your subs here" hasTooltip>
          MY SUBS
        </router-link>
        <router-link to="/newuser" v-if="!localUserExists" class="button is-link join-link">
          JOIN AETHER
        </router-link>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
  var Tooltips = require('../services/tooltips/tooltips')
  var mixins = require('../mixins/mixins')
  var globalMethods = require('../services/globals/methods')
  export default {
    name: 'a-home-header',
    mixins: [mixins.localUserMixin],
    data() {
      return {}
    },
    computed: {
      newBoardButtonDisabled(this: any) {
        if (globalMethods.IsUndefined(this.$store.state.route)) {
          return false
        }
        if (this.$store.state.route.name == 'Global>NewBoard') {
          return true
        }
        return false
      }
    },
    mounted() {
      Tooltips.Mount()
    }
  }
</script>

<style lang="scss" scoped>
  @import "../scss/bulmastyles";
  @import "../scss/globals";
  .view-header {
    height: 175px;
    display: flex;
    border-bottom: 1px solid rgba(0, 0, 0, 0.15);
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.25), 1px 0 0 0 rgba(0, 0, 0, 0.2) inset;
    background-color: $a-grey-200*0.5;
    border-left: 5px solid $a-cerulean;
    line-height: 200%;

    .side-container {
      flex: 1;
      display: flex;
    }
    .center-container {
      flex: 1;
      display: flex;
      flex-direction: column;
      .header-container {
        margin: auto;
        text-align: center;
        margin-bottom: 32px;
        &.side-layout {
          .header-text,
          .header-subtext {
            text-align: left;
            padding: 0;
          }
          .header-subtext {
            font-family: "SSP Regular"
          }
          margin-left: 50px;
        }
        .header-text {
          font-size: 200%;
        }
        .header-subtext {

          padding: 0px 50px;
        }
      }
    }
  }

  .in-button-icon {
    display: inline-block;
    width: 8px; // height: 15px;
    svg {
      height: 11px;
    }
  }

  .actions-container {
    margin: auto;
    margin-right: 50px;
    a {
      @extend %link-hover-ghost-extenders-disable;
    }
  }

  .quickstart-link:hover {
    color: $a-cerulean;
  }

  .infomark-container {
    margin-left: 3px;
    display: inline-block;
  }

  .actions-container {
    margin-bottom: 45px;
  }

  .button {
    margin-left: 5px;
  }

  .join-link {
    font-family: "SCP Bold"
  }
</style>