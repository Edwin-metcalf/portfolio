import type { DotaStatsReturn, ClashRoyaleStatsReturn } from "../../routes/games/game-stats/stats";
interface CacheEntry<T> {
    data: T | null;
    fetched: boolean;
    fetchedAt: number | null;
}
export class GameStatsCache {
    dotaData = $state<CacheEntry<DotaStatsReturn>>({
        data: null,
        fetched: false,
        fetchedAt: null,
    });
    clashRoyaleData = $state<CacheEntry<ClashRoyaleStatsReturn>>({
        data: null,
        fetched: false,
        fetchedAt: null,
    });

    loading = $state<boolean>(false)
    error = $state<string | null>(null)

    setDotaData(data : DotaStatsReturn) : void {
        this.dotaData.data = data;
        this.dotaData.fetched = true;
        this.dotaData.fetchedAt = Date.now();
    }
    setClashRoyaleData(data : ClashRoyaleStatsReturn) : void {
        this.clashRoyaleData.data = data;
        this.clashRoyaleData.fetched = true;
        this.clashRoyaleData.fetchedAt = Date.now();
    }
    isStale(fetchedAt : number | null ) : boolean {
        if (!fetchedAt) return true;
        return Date.now() - fetchedAt > 5 * 60 * 1000; //goes stale at 5 mins
    }
    
    //could use a clear data for refresh button but rather have a is stale for refresh
    //clearData() : void {
      //  this.dotaData = { data: null, fetched: false, fetchedAt: null };
      //  this.clashRoyaleData = { data: null, fetched: false, fetchedAt: null };
    //}

}

export const gameStatsCache = new GameStatsCache();
