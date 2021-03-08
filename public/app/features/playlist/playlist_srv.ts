// Libraries
import _ from 'lodash';

// Utils
import coreModule from '../../core/core_module';
import appEvents from 'app/core/app_events';

import { CoreEvents } from 'app/types';
import { getBackendSrv, locationService } from '@grafana/runtime';
import { locationUtil, urlUtil, rangeUtil } from '@grafana/data';
import { Location } from 'history';

export const queryParamsToPreserve: { [key: string]: boolean } = {
  kiosk: true,
  autofitpanels: true,
  orgId: true,
};

export class PlaylistSrv {
  private cancelPromise: any;
  private dashboards: Array<{ url: string }>;
  private index: number;
  private interval: number;
  private startUrl: string;
  private numberOfLoops = 0;
  private validPlaylistUrl: string;
  private locationListenerUnsub?: () => void;

  isPlaying: boolean;

  /** @ngInject */
  constructor(private $timeout: any) {
    this.locationUpdated = this.locationUpdated.bind(this);
  }

  next() {
    this.$timeout.cancel(this.cancelPromise);

    const playedAllDashboards = this.index > this.dashboards.length - 1;
    if (playedAllDashboards) {
      this.numberOfLoops++;

      // This does full reload of the playlist to keep memory in check due to existing leaks but at the same time
      // we do not want page to flicker after each full loop.
      if (this.numberOfLoops >= 3) {
        window.location.href = this.startUrl;
        return;
      }
      this.index = 0;
    }

    const dash = this.dashboards[this.index];
    const queryParams = locationService.getSearchObject();
    const filteredParams = _.pickBy(queryParams, (value: any, key: string) => queryParamsToPreserve[key]);
    const nextDashboardUrl = locationUtil.stripBaseFromUrl(dash.url);

    this.index++;
    this.validPlaylistUrl = nextDashboardUrl;
    this.cancelPromise = this.$timeout(() => this.next(), this.interval);

    locationService.push(nextDashboardUrl + '?' + urlUtil.toUrlParams(filteredParams));
  }

  prev() {
    this.index = Math.max(this.index - 2, 0);
    this.next();
  }

  // Detect url changes not caused by playlist srv and stop playlist
  locationUpdated(location: Location) {
    if (location.pathname !== this.validPlaylistUrl) {
      console.log('asd', location.pathname, this.validPlaylistUrl);
      this.stop();
    }
  }

  start(playlistId: number) {
    this.stop();

    this.startUrl = window.location.href;
    this.index = 0;
    this.isPlaying = true;

    // setup location tracking
    this.locationListenerUnsub = locationService.getHistory().listen(this.locationUpdated);

    appEvents.emit(CoreEvents.playlistStarted);

    return getBackendSrv()
      .get(`/api/playlists/${playlistId}`)
      .then((playlist: any) => {
        return getBackendSrv()
          .get(`/api/playlists/${playlistId}/dashboards`)
          .then((dashboards: any) => {
            this.dashboards = dashboards;
            this.interval = rangeUtil.intervalToMs(playlist.interval);
            this.next();
          });
      });
  }

  stop() {
    if (this.isPlaying) {
      const queryParams = locationService.getSearchObject();
      if (queryParams.kiosk) {
        appEvents.emit(CoreEvents.toggleKioskMode, { exit: true });
      }
    }

    this.index = 0;
    this.isPlaying = false;

    if (this.locationListenerUnsub) {
      this.locationListenerUnsub();
    }

    if (this.cancelPromise) {
      this.$timeout.cancel(this.cancelPromise);
    }

    appEvents.emit(CoreEvents.playlistStopped);
  }
}

coreModule.service('playlistSrv', PlaylistSrv);
