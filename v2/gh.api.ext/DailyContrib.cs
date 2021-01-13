using System;

namespace gh.api.ext
{
    public class DailyContrib
    {
        public int Count { get; set; }
        public Level Level { get; set; }
        public DateTimeOffset Date { get; set; }
    }
}
