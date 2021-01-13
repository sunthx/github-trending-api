using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using HtmlAgilityPack;
using Microsoft.AspNetCore.Mvc;

namespace gh.api.ext.Controllers
{
    [ApiController]
    [Route("api/gh_contrib")]
    public class GithubContributionController : ControllerBase
    {
        private string _url = "https://github.com";

        private readonly Dictionary<string, Level> _levelMap = new()
        {
            {"var(--color-calendar-graph-day-bg)", Level.L0},
            {"var(--color-calendar-graph-day-L1-bg)", Level.L1},
            {"var(--color-calendar-graph-day-L2-bg)", Level.L2},
            {"var(--color-calendar-graph-day-L3-bg)", Level.L3},
            {"var(--color-calendar-graph-day-L4-bg)", Level.L4}
        };    

        [HttpGet("{userName}")]
        public async  Task<ActionResult<ApiResult>> Get(string userName)
        {
            var result = new List<DailyContrib>();
            var requestUrl = $"{_url}/{userName}";
            var web = new HtmlWeb();
            var doc = await web.LoadFromWebAsync(requestUrl);
            var notes = doc.DocumentNode.SelectNodes("//rect[@class='day']");
            
            foreach (var htmlNode in notes)
            {
                var contrib = new DailyContrib();
                contrib.Count = htmlNode.GetAttributeValue("data-count", 0);
                
                var date = htmlNode.GetAttributeValue("data-date", string.Empty);
                if (!string.IsNullOrWhiteSpace(date))
                {
                    contrib.Date = DateTimeOffset.Parse(date);
                }

                var levelName = htmlNode.GetAttributeValue("fill", string.Empty);
                if (!string.IsNullOrWhiteSpace(levelName))
                {
                    contrib.Level = _levelMap[levelName];
                }
                
                result.Add(contrib);                                  
            }

            var apiResult = new ApiResult();
            apiResult.CreatedAt = DateTimeOffset.Now;
            apiResult.Result = result;

            return apiResult;
        }
    }
}
