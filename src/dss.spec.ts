import { Dss } from "./dss";

it('should search on google "test" and retrieve the number of results', async () =>
{
    const t = await Dss.count_event('test');
    
    expect(t).toBeGreaterThan(0);
});

it('should search on google "test test" and retrieve the number of results', async () =>
{
    const t = await Dss.count_event('test test');
    
    expect(t).toBeGreaterThan(0);
});

it('should search on google a list of value and retrieve the number of results', async () =>
{
    const t = await Dss.count_events('I like', ['cat', 'dog']);
    
    expect(t).toStrictEqual({keys: ['cat', 'dog'], values: [expect.anything(), expect.anything()]});
});

it('should search on google "io non credo nei vaccini" and retrieve the number of results', async () =>
{
    const t = await Dss.count_event('io non credo nei vaccini');
    
    expect(t).toBeGreaterThan(0);
});

it('should search on google "io amo andorea" and retrieve the number of results', async () =>
{
    const t = await Dss.count_event('io amo andorea');
    
    expect(t).toStrictEqual(0);
});