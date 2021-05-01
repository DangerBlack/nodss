import got from 'got';

const headers = {'user-agent': 'Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:87.0) Gecko/20100101 Firefox/87.0'};
const LOCALIZED_NO_RESULT_FOR_SEARCH_TERM = 'Nessun risultato trovato per i termini di ricerca';

export class Dss
{
    public static async count_event(query: string): Promise<number>
    {
        const search_key = `"${query.trim().replace(' ', '+')}"`;

        try
        {
            const result = await got(`https://www.google.it/search?q=${search_key}&gws_rd=cr,ssl&ei=odOnVsHMBcKke7jerogD`, {headers});

            if(!result.body.includes('<div id="result-stats">'))
                return 0;

            const matched = result.body.match(/<div id="result\-stats">([a-zA-Z 0-9\.]*)<nobr>/);
            const no_result = result.body.indexOf(LOCALIZED_NO_RESULT_FOR_SEARCH_TERM);

            if(no_result !== -1)
                return 0;
            
            if(!matched.length)
                return 0;

            const fields = matched[1].split(' ');
            const value = fields[fields.length > 2 ? 1 : 0].replace(/\./g,'');
            return parseInt(value, 10);
        }
        catch(error)
        {
            console.error(error);
            return 0;
        }
    }

    public static async count_events(head: string, keys: string[]): Promise<{keys: string[], values: number[]}>
    {
        const values = await Promise.all(keys.map(async (query) => {
            return await this.count_event(`${head} ${query}`);
        }));

        return {keys, values};
    }
}