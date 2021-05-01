const QuickChart = require('quickchart-js');

export async function create_pie(head_label: string, labels: string[], data: number[])
{
    const myChart = new QuickChart();

    myChart.setConfig({
        type: 'pie',
        data: {
            labels,
            datasets: [{ 
                data 
            }]
        },
        options: {
            title: {
              display: true,
              text: head_label
            },
            plugins: {
                datalabels: {
                    display: true,
                    formatter: (value, ctx) => 
                    {
                        if(value > 1000)
                            return `${value/1000}k`;
                        
                        return value;
                    },
                    color: '#000000',
                    backgroundColor: '#FFFFFF',
                    borderRadius: 3,
                },
            }
          }
    });

    return await myChart.toBinary();
}